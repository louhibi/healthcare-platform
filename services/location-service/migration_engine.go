package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"
	
	"github.com/lib/pq"
	"github.com/louhibi/healthcare-logging"
)

// SchemaMigration represents a database schema change
type SchemaMigration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// DataMigration represents data changes with dependency management
type DataMigration struct {
	Version           int
	Description       string
	RequiredSchema    int    // Must have schema version >= X
	Up                string
	Down              string
	CanRerun          bool   // Safe to run multiple times (idempotent)
	Environment       string // "all", "dev", "prod", "test"
	Tags              []string // Additional tags for filtering
}

// MigrationEngine manages both schema and data migrations
type MigrationEngine struct {
	db *sql.DB
}

// NewMigrationEngine creates a new migration engine
func NewMigrationEngine(db *sql.DB) *MigrationEngine {
	return &MigrationEngine{db: db}
}

// InitializeMigrationTables creates the migration tracking tables
func (e *MigrationEngine) InitializeMigrationTables() error {
	// Schema migrations table
	schemaQuery := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			checksum TEXT -- For integrity checking
		);
	`
	if _, err := e.db.Exec(schemaQuery); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Data migrations table
	dataQuery := `
		CREATE TABLE IF NOT EXISTS data_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			required_schema INTEGER NOT NULL,
			environment TEXT NOT NULL DEFAULT 'all',
			can_rerun BOOLEAN DEFAULT false,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			checksum TEXT,
			tags TEXT[] DEFAULT '{}'
		);
	`
	if _, err := e.db.Exec(dataQuery); err != nil {
		return fmt.Errorf("failed to create data_migrations table: %w", err)
	}

	return nil
}

// GetAppliedSchemaMigrations returns applied schema migration versions
func (e *MigrationEngine) GetAppliedSchemaMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)
	
	rows, err := e.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return applied, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// GetAppliedDataMigrations returns applied data migration versions
func (e *MigrationEngine) GetAppliedDataMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)
	
	rows, err := e.db.Query("SELECT version FROM data_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return applied, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// GetCurrentSchemaVersion returns the highest applied schema version
func (e *MigrationEngine) GetCurrentSchemaVersion() (int, error) {
	var version int
	err := e.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	return version, err
}

// RunSchemaMigrations applies pending schema migrations up to target version
func (e *MigrationEngine) RunSchemaMigrations(migrations []SchemaMigration, targetVersion int) error {
	applied, err := e.GetAppliedSchemaMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied schema migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, migration := range migrations {
		if targetVersion > 0 && migration.Version > targetVersion {
			break
		}

		if applied[migration.Version] {
			log.Printf("Schema migration %d already applied: %s", migration.Version, migration.Description)
			continue
		}

		if err := e.applySchemaMigration(migration); err != nil {
			return err
		}
	}

	return nil
}

// RunDataMigrations applies pending data migrations with dependency checking
func (e *MigrationEngine) RunDataMigrations(migrations []DataMigration, environment string, targetVersion int) error {
	currentSchemaVersion, err := e.GetCurrentSchemaVersion()
	if err != nil {
		return fmt.Errorf("failed to get current schema version: %w", err)
	}

	applied, err := e.GetAppliedDataMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied data migrations: %w", err)
	}

	// Filter migrations by environment
	var filteredMigrations []DataMigration
	for _, migration := range migrations {
		if migration.Environment == "all" || migration.Environment == environment {
			filteredMigrations = append(filteredMigrations, migration)
		}
	}

	// Sort migrations by version
	sort.Slice(filteredMigrations, func(i, j int) bool {
		return filteredMigrations[i].Version < filteredMigrations[j].Version
	})

	for _, migration := range filteredMigrations {
		if targetVersion > 0 && migration.Version > targetVersion {
			break
		}

		if applied[migration.Version] {
			if migration.CanRerun {
				log.Printf("Data migration %d already applied but can rerun: %s", migration.Version, migration.Description)
			} else {
				log.Printf("Data migration %d already applied: %s", migration.Version, migration.Description)
				continue
			}
		}

		// Check schema dependency
		if migration.RequiredSchema > currentSchemaVersion {
			return fmt.Errorf("data migration %d requires schema version %d but current version is %d", 
				migration.Version, migration.RequiredSchema, currentSchemaVersion)
		}

		if err := e.applyDataMigration(migration, environment); err != nil {
			return err
		}
	}

	return nil
}

// applySchemaMigration applies a single schema migration
func (e *MigrationEngine) applySchemaMigration(migration SchemaMigration) error {
	log.Printf("Applying schema migration %d: %s", migration.Version, migration.Description)
	
	tx, err := e.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for schema migration %d: %w", migration.Version, err)
	}
	defer tx.Rollback()

	// Execute migration with detailed logging
	logging.LogInfo("Executing schema migration SQL", "version", migration.Version, "sql", migration.Up)
	if _, err := tx.Exec(migration.Up); err != nil {
		logging.LogError("Schema migration SQL failed", "version", migration.Version, "sql", migration.Up, "error", err)
		return fmt.Errorf("failed to execute schema migration %d: %w", migration.Version, err)
	}

	// Record migration
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version, description) VALUES ($1, $2)", 
		migration.Version, migration.Description); err != nil {
		return fmt.Errorf("failed to record schema migration %d: %w", migration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit schema migration %d: %w", migration.Version, err)
	}

	log.Printf("Successfully applied schema migration %d", migration.Version)
	return nil
}

// applyDataMigration applies a single data migration
func (e *MigrationEngine) applyDataMigration(migration DataMigration, environment string) error {
	log.Printf("🔥 DEBUG: Applying data migration %d: %s (env: %s)", migration.Version, migration.Description, environment)
	
	tx, err := e.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for data migration %d: %w", migration.Version, err)
	}
	defer tx.Rollback()

	// Execute migration with detailed logging
	logging.LogInfo("Executing data migration SQL", "version", migration.Version, "sql", migration.Up)
	if _, err := tx.Exec(migration.Up); err != nil {
		logging.LogError("Data migration SQL failed", "version", migration.Version, "sql", migration.Up, "error", err)
		return fmt.Errorf("failed to execute data migration %d: %w", migration.Version, err)
	}

	// Record migration (use INSERT OR UPDATE for rerunnable migrations)
	query := `
		INSERT INTO data_migrations (version, description, required_schema, environment, can_rerun, tags) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (version) DO UPDATE SET
			applied_at = CURRENT_TIMESTAMP
	`
	if _, err := tx.Exec(query, 
		migration.Version, migration.Description, migration.RequiredSchema, 
		environment, migration.CanRerun, pq.Array(migration.Tags)); err != nil {
		return fmt.Errorf("failed to record data migration %d: %w", migration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit data migration %d: %w", migration.Version, err)
	}

	log.Printf("Successfully applied data migration %d", migration.Version)
	return nil
}

// RollbackSchemaMigration rolls back schema migrations to target version
func (e *MigrationEngine) RollbackSchemaMigration(migrations []SchemaMigration, targetVersion int) error {
	applied, err := e.GetAppliedSchemaMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied schema migrations: %w", err)
	}

	// Sort migrations by version descending for rollback
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version > migrations[j].Version
	})

	for _, migration := range migrations {
		if migration.Version <= targetVersion {
			break
		}

		if !applied[migration.Version] {
			continue
		}

		if err := e.rollbackSchemaMigration(migration); err != nil {
			return err
		}
	}

	return nil
}

// rollbackSchemaMigration rolls back a single schema migration
func (e *MigrationEngine) rollbackSchemaMigration(migration SchemaMigration) error {
	log.Printf("Rolling back schema migration %d: %s", migration.Version, migration.Description)
	
	tx, err := e.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for schema rollback %d: %w", migration.Version, err)
	}
	defer tx.Rollback()

	// Execute rollback
	if _, err := tx.Exec(migration.Down); err != nil {
		return fmt.Errorf("failed to execute schema rollback %d: %w", migration.Version, err)
	}

	// Remove migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", migration.Version); err != nil {
		return fmt.Errorf("failed to remove schema migration record %d: %w", migration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit schema rollback %d: %w", migration.Version, err)
	}

	log.Printf("Successfully rolled back schema migration %d", migration.Version)
	return nil
}

// FlushDataMigrations removes all data migration records and optionally runs down migrations
func (e *MigrationEngine) FlushDataMigrations(migrations []DataMigration, runDown bool) error {
	if runDown {
		// Sort migrations by version descending for rollback
		sort.Slice(migrations, func(i, j int) bool {
			return migrations[i].Version > migrations[j].Version
		})

		applied, err := e.GetAppliedDataMigrations()
		if err != nil {
			return fmt.Errorf("failed to get applied data migrations: %w", err)
		}

		for _, migration := range migrations {
			if !applied[migration.Version] {
				continue
			}

			log.Printf("Running down migration for data migration %d: %s", migration.Version, migration.Description)
			if _, err := e.db.Exec(migration.Down); err != nil {
				log.Printf("Warning: Failed to run down migration %d: %v", migration.Version, err)
			}
		}
	}

	// Clear all data migration records
	if _, err := e.db.Exec("DELETE FROM data_migrations"); err != nil {
		return fmt.Errorf("failed to flush data migrations: %w", err)
	}

	log.Printf("Successfully flushed all data migrations")
	return nil
}

// GetMigrationStatus returns the current status of migrations
func (e *MigrationEngine) GetMigrationStatus() (map[string]interface{}, error) {
	schemaVersion, err := e.GetCurrentSchemaVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get schema version: %w", err)
	}

	var dataCount int
	if err := e.db.QueryRow("SELECT COUNT(*) FROM data_migrations").Scan(&dataCount); err != nil {
		return nil, fmt.Errorf("failed to get data migration count: %w", err)
	}

	status := map[string]interface{}{
		"schema_version":        schemaVersion,
		"data_migrations_count": dataCount,
		"timestamp":            time.Now().Format(time.RFC3339),
	}

	return status, nil
}