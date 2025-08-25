package main

// GetSchemaMigrations returns all schema migrations in order
func GetSchemaMigrations() []SchemaMigration {
	return []SchemaMigration{
		{
			Version:     1,
			Description: "Create core location tables with indexes",
			Up: `
				CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

				-- Countries table with locale support
				CREATE TABLE countries (
					id SERIAL PRIMARY KEY,
					code VARCHAR(2) UNIQUE NOT NULL,  -- ISO 3166-1 alpha-2 code
					name_en VARCHAR(100) NOT NULL,
					name_fr VARCHAR(100),
					name_ar VARCHAR(100),
					iso_alpha3 VARCHAR(3),            -- ISO 3166-1 alpha-3 code
					numeric_code VARCHAR(3),          -- ISO 3166-1 numeric code
					region VARCHAR(50),               -- Geographic region
					subregion VARCHAR(50),            -- Geographic subregion
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- States/Provinces/Administrative divisions table with locale support
				CREATE TABLE states (
					id SERIAL PRIMARY KEY,
					country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
					code VARCHAR(10) NOT NULL,        -- State/Province code (varies by country)
					name_en VARCHAR(100) NOT NULL,
					name_fr VARCHAR(100),
					name_ar VARCHAR(100),
					type VARCHAR(20) DEFAULT 'state', -- 'state', 'province', 'region', 'territory', etc.
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(country_id, code)
				);

				-- Cities table with locale support
				CREATE TABLE cities (
					id SERIAL PRIMARY KEY,
					country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
					state_id INTEGER REFERENCES states(id) ON DELETE SET NULL,
					code VARCHAR(20),                 -- City code (optional)
					name_en VARCHAR(100) NOT NULL,
					name_fr VARCHAR(100),
					name_ar VARCHAR(100),
					latitude DECIMAL(10, 8),          -- Geographic coordinates
					longitude DECIMAL(11, 8),
					population INTEGER,               -- Population (optional)
					timezone VARCHAR(50),             -- Timezone identifier
					is_capital BOOLEAN DEFAULT FALSE, -- Is capital city
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`,
			Down: `
				DROP TABLE IF EXISTS cities CASCADE;
				DROP TABLE IF EXISTS states CASCADE;
				DROP TABLE IF EXISTS countries CASCADE;
				DROP EXTENSION IF EXISTS "uuid-ossp";
			`,
		},
		{
			Version:     2,
			Description: "Create performance indexes for location tables",
			Up: `
				-- Countries indexes
				CREATE INDEX idx_countries_code ON countries(code);
				CREATE INDEX idx_countries_name_en ON countries(name_en);
				CREATE INDEX idx_countries_name_fr ON countries(name_fr);
				CREATE INDEX idx_countries_name_ar ON countries(name_ar);
				CREATE INDEX idx_countries_active ON countries(is_active);
				CREATE INDEX idx_countries_region ON countries(region);

				-- States indexes
				CREATE INDEX idx_states_country ON states(country_id);
				CREATE INDEX idx_states_code ON states(country_id, code);
				CREATE INDEX idx_states_name_en ON states(name_en);
				CREATE INDEX idx_states_name_fr ON states(name_fr);
				CREATE INDEX idx_states_name_ar ON states(name_ar);
				CREATE INDEX idx_states_type ON states(type);
				CREATE INDEX idx_states_active ON states(is_active);

				-- Cities indexes
				CREATE INDEX idx_cities_country ON cities(country_id);
				CREATE INDEX idx_cities_state ON cities(state_id);
				CREATE INDEX idx_cities_name_en ON cities(name_en);
				CREATE INDEX idx_cities_name_fr ON cities(name_fr);
				CREATE INDEX idx_cities_name_ar ON cities(name_ar);
				CREATE INDEX idx_cities_code ON cities(code);
				CREATE INDEX idx_cities_active ON cities(is_active);
				CREATE INDEX idx_cities_capital ON cities(is_capital);
				CREATE INDEX idx_cities_population ON cities(population);
				CREATE INDEX idx_cities_coordinates ON cities(latitude, longitude);
			`,
			Down: `
				-- Drop cities indexes
				DROP INDEX IF EXISTS idx_cities_coordinates;
				DROP INDEX IF EXISTS idx_cities_population;
				DROP INDEX IF EXISTS idx_cities_capital;
				DROP INDEX IF EXISTS idx_cities_active;
				DROP INDEX IF EXISTS idx_cities_code;
				DROP INDEX IF EXISTS idx_cities_name_ar;
				DROP INDEX IF EXISTS idx_cities_name_fr;
				DROP INDEX IF EXISTS idx_cities_name_en;
				DROP INDEX IF EXISTS idx_cities_state;
				DROP INDEX IF EXISTS idx_cities_country;

				-- Drop states indexes
				DROP INDEX IF EXISTS idx_states_active;
				DROP INDEX IF EXISTS idx_states_type;
				DROP INDEX IF EXISTS idx_states_name_ar;
				DROP INDEX IF EXISTS idx_states_name_fr;
				DROP INDEX IF EXISTS idx_states_name_en;
				DROP INDEX IF EXISTS idx_states_code;
				DROP INDEX IF EXISTS idx_states_country;

				-- Drop countries indexes
				DROP INDEX IF EXISTS idx_countries_region;
				DROP INDEX IF EXISTS idx_countries_active;
				DROP INDEX IF EXISTS idx_countries_name_ar;
				DROP INDEX IF EXISTS idx_countries_name_fr;
				DROP INDEX IF EXISTS idx_countries_name_en;
				DROP INDEX IF EXISTS idx_countries_code;
			`,
		},
		{
			Version:     3,
			Description: "Create updated_at triggers for timestamp management",
			Up: `
				-- Updated timestamp trigger function
				CREATE OR REPLACE FUNCTION update_updated_at_column()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.updated_at = CURRENT_TIMESTAMP;
					RETURN NEW;
				END;
				$$ language 'plpgsql';

				-- Apply triggers to all location tables
				CREATE TRIGGER update_countries_updated_at BEFORE UPDATE ON countries
					FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

				CREATE TRIGGER update_states_updated_at BEFORE UPDATE ON states
					FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

				CREATE TRIGGER update_cities_updated_at BEFORE UPDATE ON cities
					FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
			`,
			Down: `
				DROP TRIGGER IF EXISTS update_cities_updated_at ON cities;
				DROP TRIGGER IF EXISTS update_states_updated_at ON states;
				DROP TRIGGER IF EXISTS update_countries_updated_at ON countries;
				DROP FUNCTION IF EXISTS update_updated_at_column();
			`,
		},
		{
			Version:     4,
			Description: "Add search optimization features for location lookup",
			Up: `
				-- Add full-text search capabilities
				ALTER TABLE countries 
				ADD COLUMN search_vector tsvector;

				ALTER TABLE states
				ADD COLUMN search_vector tsvector;

				ALTER TABLE cities
				ADD COLUMN search_vector tsvector;

				-- Create search vector update triggers
				CREATE OR REPLACE FUNCTION update_countries_search_vector()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.search_vector = 
						setweight(to_tsvector('simple', COALESCE(NEW.name_en, '')), 'A') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_fr, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_ar, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.code, '')), 'A');
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;

				CREATE OR REPLACE FUNCTION update_states_search_vector()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.search_vector = 
						setweight(to_tsvector('simple', COALESCE(NEW.name_en, '')), 'A') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_fr, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_ar, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.code, '')), 'A');
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;

				CREATE OR REPLACE FUNCTION update_cities_search_vector()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.search_vector = 
						setweight(to_tsvector('simple', COALESCE(NEW.name_en, '')), 'A') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_fr, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_ar, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.code, '')), 'A');
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;

				-- Create search triggers
				CREATE TRIGGER countries_search_vector_trigger
					BEFORE INSERT OR UPDATE ON countries
					FOR EACH ROW EXECUTE FUNCTION update_countries_search_vector();

				CREATE TRIGGER states_search_vector_trigger
					BEFORE INSERT OR UPDATE ON states
					FOR EACH ROW EXECUTE FUNCTION update_states_search_vector();

				CREATE TRIGGER cities_search_vector_trigger
					BEFORE INSERT OR UPDATE ON cities
					FOR EACH ROW EXECUTE FUNCTION update_cities_search_vector();

				-- Create search indexes
				CREATE INDEX idx_countries_search ON countries USING GIN(search_vector);
				CREATE INDEX idx_states_search ON states USING GIN(search_vector);  
				CREATE INDEX idx_cities_search ON cities USING GIN(search_vector);
			`,
			Down: `
				-- Drop search indexes
				DROP INDEX IF EXISTS idx_cities_search;
				DROP INDEX IF EXISTS idx_states_search;
				DROP INDEX IF EXISTS idx_countries_search;

				-- Drop search triggers
				DROP TRIGGER IF EXISTS cities_search_vector_trigger ON cities;
				DROP TRIGGER IF EXISTS states_search_vector_trigger ON states;
				DROP TRIGGER IF EXISTS countries_search_vector_trigger ON countries;

				-- Drop search functions
				DROP FUNCTION IF EXISTS update_cities_search_vector();
				DROP FUNCTION IF EXISTS update_states_search_vector();
				DROP FUNCTION IF EXISTS update_countries_search_vector();

				-- Drop search vector columns
				ALTER TABLE cities DROP COLUMN IF EXISTS search_vector;
				ALTER TABLE states DROP COLUMN IF EXISTS search_vector;
				ALTER TABLE countries DROP COLUMN IF EXISTS search_vector;
			`,
		},
		{
			Version:     5,
			Description: "Create nationalities table linked to countries",
			Up: `
				-- Nationalities table with locale support
				CREATE TABLE nationalities (
					id SERIAL PRIMARY KEY,
					country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
					name_en VARCHAR(100) NOT NULL,
					name_fr VARCHAR(100),
					name_ar VARCHAR(100),
					code VARCHAR(10),                 -- Optional nationality code (e.g., "CA", "US", "MA", "FR")
					is_primary BOOLEAN DEFAULT FALSE, -- Is this the primary nationality for the country
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(country_id, code)
				);

				-- Indexes for nationalities
				CREATE INDEX idx_nationalities_country ON nationalities(country_id);
				CREATE INDEX idx_nationalities_name_en ON nationalities(name_en);
				CREATE INDEX idx_nationalities_name_fr ON nationalities(name_fr);
				CREATE INDEX idx_nationalities_name_ar ON nationalities(name_ar);
				CREATE INDEX idx_nationalities_code ON nationalities(code);
				CREATE INDEX idx_nationalities_primary ON nationalities(is_primary);
				CREATE INDEX idx_nationalities_active ON nationalities(is_active);

				-- Add updated_at trigger for nationalities
				CREATE TRIGGER update_nationalities_updated_at BEFORE UPDATE ON nationalities
					FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

				-- Add search vector column and trigger for nationalities
				ALTER TABLE nationalities 
				ADD COLUMN search_vector tsvector;

				CREATE OR REPLACE FUNCTION update_nationalities_search_vector()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.search_vector = 
						setweight(to_tsvector('simple', COALESCE(NEW.name_en, '')), 'A') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_fr, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.name_ar, '')), 'B') ||
						setweight(to_tsvector('simple', COALESCE(NEW.code, '')), 'A');
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;

				CREATE TRIGGER nationalities_search_vector_trigger
					BEFORE INSERT OR UPDATE ON nationalities
					FOR EACH ROW EXECUTE FUNCTION update_nationalities_search_vector();

				CREATE INDEX idx_nationalities_search ON nationalities USING GIN(search_vector);
			`,
			Down: `
				-- Drop nationalities search components
				DROP INDEX IF EXISTS idx_nationalities_search;
				DROP TRIGGER IF EXISTS nationalities_search_vector_trigger ON nationalities;
				DROP FUNCTION IF EXISTS update_nationalities_search_vector();

				-- Drop nationalities triggers
				DROP TRIGGER IF EXISTS update_nationalities_updated_at ON nationalities;

				-- Drop nationalities indexes
				DROP INDEX IF EXISTS idx_nationalities_active;
				DROP INDEX IF EXISTS idx_nationalities_primary;
				DROP INDEX IF EXISTS idx_nationalities_code;
				DROP INDEX IF EXISTS idx_nationalities_name_ar;
				DROP INDEX IF EXISTS idx_nationalities_name_fr;
				DROP INDEX IF EXISTS idx_nationalities_name_en;
				DROP INDEX IF EXISTS idx_nationalities_country;

				-- Drop nationalities table
				DROP TABLE IF EXISTS nationalities CASCADE;
			`,
		},
		{
			Version:     4,
			Description: "Create insurance tables with multi-locale support",
			Up: `
				-- Insurance Types table with locale support
				CREATE TABLE insurance_types (
					id SERIAL PRIMARY KEY,
					country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
					code VARCHAR(50) NOT NULL,           -- Type code (e.g., "public", "private", "other")
					name_en VARCHAR(100) NOT NULL,       -- English name
					name_fr VARCHAR(100),                -- French name
					name_ar VARCHAR(100),                -- Arabic name
					is_default BOOLEAN DEFAULT FALSE,    -- Is this the default type for the country
					sort_order INTEGER DEFAULT 0,       -- Display order
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(country_id, code)
				);

				-- Insurance Providers table with locale support
				CREATE TABLE insurance_providers (
					id SERIAL PRIMARY KEY,
					insurance_type_id INTEGER NOT NULL REFERENCES insurance_types(id) ON DELETE CASCADE,
					code VARCHAR(50) NOT NULL,           -- Provider code (e.g., "ohip", "medicare", "other")
					name_en VARCHAR(100) NOT NULL,       -- English name
					name_fr VARCHAR(100),                -- French name
					name_ar VARCHAR(100),                -- Arabic name
					is_default BOOLEAN DEFAULT FALSE,    -- Is this the default provider for the type
					sort_order INTEGER DEFAULT 0,       -- Display order
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(insurance_type_id, code)
				);

				-- Add indexes for insurance types
				CREATE INDEX idx_insurance_types_country ON insurance_types(country_id);
				CREATE INDEX idx_insurance_types_code ON insurance_types(code);
				CREATE INDEX idx_insurance_types_active ON insurance_types(is_active);
				CREATE INDEX idx_insurance_types_default ON insurance_types(is_default);
				CREATE INDEX idx_insurance_types_sort ON insurance_types(sort_order);
				CREATE INDEX idx_insurance_types_name_en ON insurance_types(name_en);
				CREATE INDEX idx_insurance_types_name_fr ON insurance_types(name_fr);
				CREATE INDEX idx_insurance_types_name_ar ON insurance_types(name_ar);

				-- Add indexes for insurance providers
				CREATE INDEX idx_insurance_providers_type ON insurance_providers(insurance_type_id);
				CREATE INDEX idx_insurance_providers_code ON insurance_providers(code);
				CREATE INDEX idx_insurance_providers_active ON insurance_providers(is_active);
				CREATE INDEX idx_insurance_providers_default ON insurance_providers(is_default);
				CREATE INDEX idx_insurance_providers_sort ON insurance_providers(sort_order);
				CREATE INDEX idx_insurance_providers_name_en ON insurance_providers(name_en);
				CREATE INDEX idx_insurance_providers_name_fr ON insurance_providers(name_fr);
				CREATE INDEX idx_insurance_providers_name_ar ON insurance_providers(name_ar);

				-- Add updated_at trigger for insurance_types
				CREATE TRIGGER update_insurance_types_updated_at BEFORE UPDATE ON insurance_types
					FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

				-- Add updated_at trigger for insurance_providers
				CREATE TRIGGER update_insurance_providers_updated_at BEFORE UPDATE ON insurance_providers
					FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
			`,
			Down: `
				-- Drop insurance providers triggers
				DROP TRIGGER IF EXISTS update_insurance_providers_updated_at ON insurance_providers;

				-- Drop insurance types triggers
				DROP TRIGGER IF EXISTS update_insurance_types_updated_at ON insurance_types;

				-- Drop insurance providers indexes
				DROP INDEX IF EXISTS idx_insurance_providers_name_ar;
				DROP INDEX IF EXISTS idx_insurance_providers_name_fr;
				DROP INDEX IF EXISTS idx_insurance_providers_name_en;
				DROP INDEX IF EXISTS idx_insurance_providers_sort;
				DROP INDEX IF EXISTS idx_insurance_providers_default;
				DROP INDEX IF EXISTS idx_insurance_providers_active;
				DROP INDEX IF EXISTS idx_insurance_providers_code;
				DROP INDEX IF EXISTS idx_insurance_providers_type;

				-- Drop insurance types indexes
				DROP INDEX IF EXISTS idx_insurance_types_name_ar;
				DROP INDEX IF EXISTS idx_insurance_types_name_fr;
				DROP INDEX IF EXISTS idx_insurance_types_name_en;
				DROP INDEX IF EXISTS idx_insurance_types_sort;
				DROP INDEX IF EXISTS idx_insurance_types_default;
				DROP INDEX IF EXISTS idx_insurance_types_active;
				DROP INDEX IF EXISTS idx_insurance_types_code;
				DROP INDEX IF EXISTS idx_insurance_types_country;

				-- Drop insurance tables
				DROP TABLE IF EXISTS insurance_providers CASCADE;
				DROP TABLE IF EXISTS insurance_types CASCADE;
			`,
		},
	}
}