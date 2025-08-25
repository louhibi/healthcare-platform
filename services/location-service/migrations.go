package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// GetMigrations returns all migrations in order
func GetMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "Create location tables with schema",
			Up: `
				CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

				-- Countries table
				CREATE TABLE countries (
					id SERIAL PRIMARY KEY,
					code VARCHAR(2) UNIQUE NOT NULL,  -- ISO 3166-1 alpha-2 code
					name VARCHAR(100) NOT NULL,
					iso_alpha3 VARCHAR(3),            -- ISO 3166-1 alpha-3 code
					numeric_code VARCHAR(3),          -- ISO 3166-1 numeric code
					region VARCHAR(50),               -- Geographic region
					subregion VARCHAR(50),            -- Geographic subregion
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- States/Provinces/Administrative divisions table
				CREATE TABLE states (
					id SERIAL PRIMARY KEY,
					country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
					code VARCHAR(10) NOT NULL,        -- State/Province code (varies by country)
					name VARCHAR(100) NOT NULL,
					type VARCHAR(20) DEFAULT 'state', -- 'state', 'province', 'region', 'territory', etc.
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(country_id, code)
				);

				-- Cities table
				CREATE TABLE cities (
					id SERIAL PRIMARY KEY,
					country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
					state_id INTEGER REFERENCES states(id) ON DELETE SET NULL,
					code VARCHAR(20),                 -- City code (optional)
					name VARCHAR(100) NOT NULL,
					latitude DECIMAL(10, 8),          -- Geographic coordinates
					longitude DECIMAL(11, 8),
					population INTEGER,               -- Population (optional)
					timezone VARCHAR(50),             -- Timezone identifier
					is_capital BOOLEAN DEFAULT FALSE, -- Is capital city
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Indexes for performance
				CREATE INDEX idx_countries_code ON countries(code);
				CREATE INDEX idx_countries_name ON countries(name);
				CREATE INDEX idx_countries_active ON countries(is_active);

				CREATE INDEX idx_states_country ON states(country_id);
				CREATE INDEX idx_states_code ON states(country_id, code);
				CREATE INDEX idx_states_name ON states(name);
				CREATE INDEX idx_states_active ON states(is_active);

				CREATE INDEX idx_cities_country ON cities(country_id);
				CREATE INDEX idx_cities_state ON cities(state_id);
				CREATE INDEX idx_cities_name ON cities(name);
				CREATE INDEX idx_cities_code ON cities(code);
				CREATE INDEX idx_cities_active ON cities(is_active);
				CREATE INDEX idx_cities_capital ON cities(is_capital);

				-- Updated timestamp triggers
				CREATE OR REPLACE FUNCTION update_updated_at_column()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.updated_at = CURRENT_TIMESTAMP;
					RETURN NEW;
				END;
				$$ language 'plpgsql';

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
				DROP TABLE IF EXISTS cities CASCADE;
				DROP TABLE IF EXISTS states CASCADE;
				DROP TABLE IF EXISTS countries CASCADE;
				DROP EXTENSION IF EXISTS "uuid-ossp";
			`,
		},
		{
			Version:     2,
			Description: "Load location fixtures for healthcare platform countries",
			Up: `
				-- Countries (Healthcare platform supported countries)
				INSERT INTO countries (code, name, iso_alpha3, numeric_code, region, subregion) VALUES
				('CA', 'Canada', 'CAN', '124', 'Americas', 'Northern America'),
				('US', 'United States', 'USA', '840', 'Americas', 'Northern America'),
				('MA', 'Morocco', 'MAR', '504', 'Africa', 'Northern Africa'),
				('FR', 'France', 'FRA', '250', 'Europe', 'Western Europe');

				-- Canadian Provinces and Territories
				INSERT INTO states (country_id, code, name, type) VALUES
				((SELECT id FROM countries WHERE code = 'CA'), 'ON', 'Ontario', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'QC', 'Quebec', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'BC', 'British Columbia', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'AB', 'Alberta', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'MB', 'Manitoba', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'SK', 'Saskatchewan', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NS', 'Nova Scotia', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NB', 'New Brunswick', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NL', 'Newfoundland and Labrador', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'PE', 'Prince Edward Island', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'YT', 'Yukon', 'territory'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NT', 'Northwest Territories', 'territory'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NU', 'Nunavut', 'territory');

				-- US States (major states for healthcare platform)
				INSERT INTO states (country_id, code, name, type) VALUES
				((SELECT id FROM countries WHERE code = 'US'), 'AL', 'Alabama', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'AK', 'Alaska', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'CA', 'California', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'FL', 'Florida', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'GA', 'Georgia', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'IL', 'Illinois', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MA', 'Massachusetts', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MI', 'Michigan', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NY', 'New York', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NC', 'North Carolina', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'OH', 'Ohio', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'PA', 'Pennsylvania', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'TX', 'Texas', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'WA', 'Washington', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'DC', 'District of Columbia', 'district');

				-- Moroccan Regions
				INSERT INTO states (country_id, code, name, type) VALUES
				((SELECT id FROM countries WHERE code = 'MA'), '01', 'Tanger-Tétouan-Al Hoceima', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '02', 'L''Oriental', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '03', 'Fès-Meknès', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '04', 'Rabat-Salé-Kénitra', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '05', 'Béni Mellal-Khénifra', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '06', 'Casablanca-Settat', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '07', 'Marrakech-Safi', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '08', 'Drâa-Tafilalet', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '09', 'Souss-Massa', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '10', 'Guelmim-Oued Noun', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '11', 'Laâyoune-Sakia El Hamra', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '12', 'Dakhla-Oued Ed-Dahab', 'region');

				-- French Regions
				INSERT INTO states (country_id, code, name, type) VALUES
				((SELECT id FROM countries WHERE code = 'FR'), 'ARA', 'Auvergne-Rhône-Alpes', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'BFC', 'Bourgogne-Franche-Comté', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'BRE', 'Bretagne', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'CVL', 'Centre-Val de Loire', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'COR', 'Corse', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'GES', 'Grand Est', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'HDF', 'Hauts-de-France', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'IDF', 'Île-de-France', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'NOR', 'Normandie', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'NAQ', 'Nouvelle-Aquitaine', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'OCC', 'Occitanie', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'PDL', 'Pays de la Loire', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'PAC', 'Provence-Alpes-Côte d''Azur', 'region');
			`,
			Down: `
				DELETE FROM states;
				DELETE FROM countries;
			`,
		},
		{
			Version:     3,
			Description: "Load city fixtures for healthcare platform",
			Up: `
				-- Major Canadian Cities
				INSERT INTO cities (country_id, state_id, code, name, latitude, longitude, population, timezone, is_capital) VALUES
				-- Ontario
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'toronto', 'Toronto', 43.6532, -79.3832, 2930000, 'America/Toronto', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'ottawa', 'Ottawa', 45.4215, -75.6972, 994000, 'America/Toronto', TRUE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'hamilton', 'Hamilton', 43.2557, -79.8711, 537000, 'America/Toronto', FALSE),
				-- Quebec
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'QC'), 'montreal', 'Montreal', 45.5017, -73.5673, 1780000, 'America/Toronto', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'QC'), 'quebec', 'Quebec City', 46.8139, -71.2080, 542000, 'America/Toronto', FALSE),
				-- British Columbia
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'BC'), 'vancouver', 'Vancouver', 49.2827, -123.1207, 675000, 'America/Vancouver', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'BC'), 'victoria', 'Victoria', 48.4284, -123.3656, 91000, 'America/Vancouver', FALSE),
				-- Alberta
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'AB'), 'calgary', 'Calgary', 51.0447, -114.0719, 1306000, 'America/Edmonton', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'AB'), 'edmonton', 'Edmonton', 53.5461, -113.4938, 981000, 'America/Edmonton', FALSE);

				-- Major US Cities
				INSERT INTO cities (country_id, state_id, code, name, latitude, longitude, population, timezone, is_capital) VALUES
				-- New York
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'NY'), 'nyc', 'New York City', 40.7128, -74.0060, 8336000, 'America/New_York', FALSE),
				-- California
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'CA'), 'la', 'Los Angeles', 34.0522, -118.2437, 3980000, 'America/Los_Angeles', FALSE),
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'CA'), 'sf', 'San Francisco', 37.7749, -122.4194, 874000, 'America/Los_Angeles', FALSE),
				-- Texas
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'TX'), 'houston', 'Houston', 29.7604, -95.3698, 2320000, 'America/Chicago', FALSE),
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'TX'), 'dallas', 'Dallas', 32.7767, -96.7970, 1343000, 'America/Chicago', FALSE),
				-- Illinois
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'IL'), 'chicago', 'Chicago', 41.8781, -87.6298, 2716000, 'America/Chicago', FALSE);

				-- Major Moroccan Cities
				INSERT INTO cities (country_id, state_id, code, name, latitude, longitude, population, timezone, is_capital) VALUES
				-- Casablanca-Settat
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'casa', 'Casablanca', 33.5731, -7.5898, 3360000, 'Africa/Casablanca', FALSE),
				-- Rabat-Salé-Kénitra
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'rabat', 'Rabat', 34.0209, -6.8416, 580000, 'Africa/Casablanca', TRUE),
				-- Fès-Meknès
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'fes', 'Fès', 34.0181, -5.0078, 1112000, 'Africa/Casablanca', FALSE),
				-- Marrakech-Safi
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'marrakech', 'Marrakech', 31.6295, -7.9811, 928000, 'Africa/Casablanca', FALSE);

				-- Major French Cities
				INSERT INTO cities (country_id, state_id, code, name, latitude, longitude, population, timezone, is_capital) VALUES
				-- Île-de-France
				((SELECT id FROM countries WHERE code = 'FR'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'FR') AND code = 'IDF'), 'paris', 'Paris', 48.8566, 2.3522, 2161000, 'Europe/Paris', TRUE),
				-- Auvergne-Rhône-Alpes
				((SELECT id FROM countries WHERE code = 'FR'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'FR') AND code = 'ARA'), 'lyon', 'Lyon', 45.7640, 4.8357, 518000, 'Europe/Paris', FALSE),
				-- Provence-Alpes-Côte d'Azur
				((SELECT id FROM countries WHERE code = 'FR'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'FR') AND code = 'PAC'), 'marseille', 'Marseille', 43.2965, 5.3698, 870000, 'Europe/Paris', FALSE),
				-- Occitanie
				((SELECT id FROM countries WHERE code = 'FR'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'FR') AND code = 'OCC'), 'toulouse', 'Toulouse', 43.6047, 1.4442, 479000, 'Europe/Paris', FALSE);
			`,
			Down: `
				DELETE FROM cities;
			`,
		},
		{
			Version:     4,
			Description: "Add comprehensive Moroccan cities data",
			Up: `
				-- Comprehensive Moroccan Cities by Region
				INSERT INTO cities (country_id, state_id, code, name, latitude, longitude, population, timezone, is_capital) VALUES
				
				-- Tanger-Tétouan-Al Hoceima (Region 01)
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'tangier', 'Tanger', 35.7595, -5.8340, 947952, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'tetouan', 'Tétouan', 35.5889, -5.3626, 380787, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'alhoceima', 'Al Hoceïma', 35.2517, -3.9317, 56716, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'chefchaouen', 'Chefchaouen', 35.1688, -5.2636, 42786, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'larache', 'Larache', 35.1932, -6.1561, 125008, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'ksar_elkebir', 'Ksar el-Kebir', 35.0056, -5.9018, 126617, 'Africa/Casablanca', FALSE),

				-- L'Oriental (Region 02)  
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'oujda', 'Oujda', 34.6814, -1.9086, 494252, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'nador', 'Nador', 35.1681, -2.9287, 178540, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'berkane', 'Berkane', 34.9183, -2.3217, 109237, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'taourirt', 'Taourirt', 34.4092, -2.8953, 82518, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'jerada', 'Jerada', 34.3117, -2.1631, 43506, 'Africa/Casablanca', FALSE),

				-- Fès-Meknès (Region 03) - Add more cities to existing Fès
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'meknes', 'Meknès', 33.8935, -5.5473, 632079, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'taza', 'Taza', 34.2133, -4.0100, 148456, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'sefrou', 'Sefrou', 33.8275, -4.8378, 79887, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'errachidia', 'Errachidia', 31.9314, -4.4241, 92374, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'ifrane', 'Ifrane', 33.5228, -5.1106, 14659, 'Africa/Casablanca', FALSE),

				-- Rabat-Salé-Kénitra (Region 04) - Add more cities to existing Rabat
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'sale', 'Salé', 34.0531, -6.7986, 890403, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'kenitra', 'Kénitra', 34.2610, -6.5802, 431282, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'temara', 'Témara', 33.9287, -6.9063, 344924, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'skhirate', 'Skhirate', 33.8561, -7.0374, 38026, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'khemisset', 'Khemisset', 33.8244, -6.0668, 127917, 'Africa/Casablanca', FALSE),

				-- Béni Mellal-Khénifra (Region 05)
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'beni_mellal', 'Béni Mellal', 32.3372, -6.3498, 192676, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'khenifra', 'Khénifra', 32.9353, -5.6681, 86716, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'fkih_ben_salah', 'Fkih Ben Salah', 32.5017, -6.6914, 78840, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'kasba_tadla', 'Kasba Tadla', 32.5975, -6.2656, 44598, 'Africa/Casablanca', FALSE),

				-- Casablanca-Settat (Region 06) - Add more cities to existing Casablanca
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'mohammedia', 'Mohammedia', 33.6866, -7.3837, 322286, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'settat', 'Settat', 33.0008, -7.6164, 142250, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'berrechid', 'Berrechid', 33.2650, -7.5869, 136634, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'benslimane', 'Benslimane', 33.6133, -7.1211, 46522, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'nouaceur', 'Nouaceur', 33.3983, -7.5397, 42817, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'el_jadida', 'El Jadida', 33.2316, -8.5007, 196789, 'Africa/Casablanca', FALSE),

				-- Marrakech-Safi (Region 07) - Add more cities to existing Marrakech
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'safi', 'Safi', 32.2994, -9.2372, 308508, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'essaouira', 'Essaouira', 31.5085, -9.7595, 77966, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'kelaa_des_sraghna', 'Kelâa des Sraghna', 32.0544, -7.4059, 98943, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'youssoufia', 'Youssoufia', 32.2465, -8.5286, 67628, 'Africa/Casablanca', FALSE),

				-- Drâa-Tafilalet (Region 08)
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'ouarzazate', 'Ouarzazate', 30.9189, -6.8939, 71067, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'zagora', 'Zagora', 30.3273, -5.8372, 40069, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'tinghir', 'Tinghir', 31.5147, -5.5333, 43632, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'rissani', 'Rissani', 31.2794, -4.2594, 20469, 'Africa/Casablanca', FALSE),

				-- Souss-Massa (Region 09)
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'agadir', 'Agadir', 30.4278, -9.5981, 421844, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'inezgane', 'Inezgane', 30.3556, -9.5361, 101243, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'taroudant', 'Taroudant', 30.4731, -8.8779, 92845, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'tiznit', 'Tiznit', 29.6978, -9.7316, 74699, 'Africa/Casablanca', FALSE),

				-- Guelmim-Oued Noun (Region 10)
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '10'), 'guelmim', 'Guelmim', 28.9870, -10.0574, 118318, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '10'), 'tan_tan', 'Tan-Tan', 28.4378, -11.103, 73209, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '10'), 'sidi_ifni', 'Sidi Ifni', 29.3792, -10.1731, 20051, 'Africa/Casablanca', FALSE),

				-- Laâyoune-Sakia El Hamra (Region 11)
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '11'), 'laayoune', 'Laâyoune', 27.1536, -13.1994, 271480, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '11'), 'boujdour', 'Boujdour', 26.1265, -14.4816, 46129, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '11'), 'smara', 'Smara', 26.7316, -11.6719, 57035, 'Africa/Casablanca', FALSE),

				-- Dakhla-Oued Ed-Dahab (Region 12)
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '12'), 'dakhla', 'Dakhla', 23.6848, -15.9501, 106277, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '12'), 'aousserd', 'Aousserd', 22.5667, -14.0833, 7609, 'Africa/Casablanca', FALSE);
			`,
			Down: `
				-- Remove only the newly added Moroccan cities
				DELETE FROM cities WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') 
				AND code NOT IN ('casa', 'rabat', 'fes', 'marrakech');
			`,
		},
	}
}

// CreateMigrationsTable creates the schema_migrations table to track applied migrations
func CreateMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

// GetAppliedMigrations returns a map of applied migration versions
func GetAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	applied := make(map[int]bool)
	
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
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

// RunMigrations applies all pending migrations
func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := CreateMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrations := GetMigrations()
	
	for _, migration := range migrations {
		if applied[migration.Version] {
			log.Printf("Migration %d already applied: %s", migration.Version, migration.Description)
			continue
		}

		log.Printf("Applying migration %d: %s", migration.Version, migration.Description)
		
		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
		}

		// Execute migration
		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version, description) VALUES ($1, $2)", 
			migration.Version, migration.Description); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		log.Printf("Successfully applied migration %d", migration.Version)
	}

	return nil
}