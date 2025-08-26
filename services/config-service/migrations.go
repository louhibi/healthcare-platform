package main

import (
    "database/sql"
    "fmt"
    "log"
)

type Migration struct {
    Version int
    Description string
    Up string
    Down string
}

func getMigrations() []Migration {
    return []Migration{
        {
            Version: 1,
            Description: "Create config_settings and feature_flags tables",
            Up: `
                CREATE TABLE IF NOT EXISTS config_settings (
                    id SERIAL PRIMARY KEY,
                    key VARCHAR(100) UNIQUE NOT NULL,
                    value TEXT NOT NULL,
                    is_public BOOLEAN DEFAULT false,
                    description TEXT,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                );
                CREATE TABLE IF NOT EXISTS feature_flags (
                    id SERIAL PRIMARY KEY,
                    name VARCHAR(100) UNIQUE NOT NULL,
                    enabled BOOLEAN DEFAULT false,
                    description TEXT,
                    is_public BOOLEAN DEFAULT true,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                );
                -- Seed default bootstrap config
                INSERT INTO config_settings (key, value, is_public, description) VALUES
                    ('app.name', 'Healthcare Platform', true, 'Application display name'),
                    ('app.version', '0.1.0', true, 'Application version'),
                    ('auth.password.min_length', '8', false, 'Minimum password length')
                ON CONFLICT (key) DO NOTHING;
                INSERT INTO feature_flags (name, enabled, description, is_public) VALUES
                    ('appointments.enabled', true, 'Enable appointment booking', true),
                    ('registration.enabled', true, 'Allow public self-registration', true)
                ON CONFLICT (name) DO NOTHING;
            `,
            Down: `DROP TABLE IF EXISTS feature_flags; DROP TABLE IF EXISTS config_settings;`,
        },
    }
}

func RunMigrationsWithVersioning(db *sql.DB) error {
    if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version INT PRIMARY KEY, applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`); err != nil { return err }
    migrations := getMigrations()
    for _, m := range migrations {
        var exists bool
        if err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)", m.Version).Scan(&exists); err != nil { return err }
        if exists { continue }
        log.Printf("Applying migration %d: %s", m.Version, m.Description)
        tx, err := db.Begin(); if err != nil { return err }
        if _, err := tx.Exec(m.Up); err != nil { tx.Rollback(); return fmt.Errorf("migration %d failed: %w", m.Version, err) }
        if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", m.Version); err != nil { tx.Rollback(); return err }
        if err := tx.Commit(); err != nil { return err }
    }
    return nil
}
