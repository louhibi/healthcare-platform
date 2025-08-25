package main

// GetDataMigrations returns all data migrations in order
func GetDataMigrations() []DataMigration {
	return []DataMigration{
		{
			Version:        1,
			Description:    "Load core healthcare platform countries",
			RequiredSchema: 1,
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"countries", "core"},
			Up: `
				-- Core healthcare platform countries with full locale support
				INSERT INTO countries (code, name_en, name_fr, name_ar, iso_alpha3, numeric_code, region, subregion) VALUES
				('CA', 'Canada', 'Canada', 'كندا', 'CAN', '124', 'Americas', 'Northern America'),
				('US', 'United States', 'États-Unis', 'الولايات المتحدة', 'USA', '840', 'Americas', 'Northern America'),
				('MA', 'Morocco', 'Maroc', 'المغرب', 'MAR', '504', 'Africa', 'Northern Africa'),
				('FR', 'France', 'France', 'فرنسا', 'FRA', '250', 'Europe', 'Western Europe')
				ON CONFLICT (code) DO UPDATE SET
					name_en = EXCLUDED.name_en,
					name_fr = EXCLUDED.name_fr,
					name_ar = EXCLUDED.name_ar,
					iso_alpha3 = EXCLUDED.iso_alpha3,
					numeric_code = EXCLUDED.numeric_code,
					region = EXCLUDED.region,
					subregion = EXCLUDED.subregion,
					updated_at = CURRENT_TIMESTAMP;
			`,
			Down: `
				DELETE FROM countries WHERE code IN ('CA', 'US', 'MA', 'FR');
			`,
		},
		{
			Version:        2,
			Description:    "Load Canadian provinces and territories",
			RequiredSchema: 1,
			Environment:    "all", 
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"states", "canada"},
			Up: `
				-- Canadian Provinces and Territories with full locale support
				INSERT INTO states (country_id, code, name_en, name_fr, name_ar, type) VALUES
				((SELECT id FROM countries WHERE code = 'CA'), 'ON', 'Ontario', 'Ontario', 'أونتاريو', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'QC', 'Quebec', 'Québec', 'كيبيك', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'BC', 'British Columbia', 'Colombie-Britannique', 'كولومبيا البريطانية', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'AB', 'Alberta', 'Alberta', 'ألبرتا', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'MB', 'Manitoba', 'Manitoba', 'مانيتوبا', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'SK', 'Saskatchewan', 'Saskatchewan', 'ساسكاتشوان', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NS', 'Nova Scotia', 'Nouvelle-Écosse', 'نوفا سكوتيا', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NB', 'New Brunswick', 'Nouveau-Brunswick', 'نيو برونزويك', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NL', 'Newfoundland and Labrador', 'Terre-Neuve-et-Labrador', 'نيوفاوندلاند ولابرادور', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'PE', 'Prince Edward Island', 'Île-du-Prince-Édouard', 'جزيرة الأمير إدوارد', 'province'),
				((SELECT id FROM countries WHERE code = 'CA'), 'YT', 'Yukon', 'Yukon', 'يوكون', 'territory'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NT', 'Northwest Territories', 'Territoires du Nord-Ouest', 'الأقاليم الشمالية الغربية', 'territory'),
				((SELECT id FROM countries WHERE code = 'CA'), 'NU', 'Nunavut', 'Nunavut', 'نونافوت', 'territory')
				-- Insert only, no conflicts expected
			`,
			Down: `
				DELETE FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA');
			`,
		},
		{
			Version:        3,
			Description:    "Load US states and territories",
			RequiredSchema: 1,
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"states", "usa"},
			Up: `
				-- US States and Territories with full locale support
				INSERT INTO states (country_id, code, name_en, name_fr, name_ar, type) VALUES
				((SELECT id FROM countries WHERE code = 'US'), 'AL', 'Alabama', 'Alabama', 'ألاباما', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'AK', 'Alaska', 'Alaska', 'ألاسكا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'AZ', 'Arizona', 'Arizona', 'أريزونا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'AR', 'Arkansas', 'Arkansas', 'أركنساس', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'CA', 'California', 'Californie', 'كاليفورنيا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'CO', 'Colorado', 'Colorado', 'كولورادو', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'CT', 'Connecticut', 'Connecticut', 'كونيتيكت', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'DE', 'Delaware', 'Delaware', 'ديلاوير', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'FL', 'Florida', 'Floride', 'فلوريدا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'GA', 'Georgia', 'Géorgie', 'جورجيا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'HI', 'Hawaii', 'Hawaï', 'هاواي', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'ID', 'Idaho', 'Idaho', 'أيداهو', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'IL', 'Illinois', 'Illinois', 'إلينوي', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'IN', 'Indiana', 'Indiana', 'إنديانا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'IA', 'Iowa', 'Iowa', 'آيوا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'KS', 'Kansas', 'Kansas', 'كانساس', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'KY', 'Kentucky', 'Kentucky', 'كنتاكي', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'LA', 'Louisiana', 'Louisiane', 'لويزيانا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'ME', 'Maine', 'Maine', 'مين', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MD', 'Maryland', 'Maryland', 'ماريلاند', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MA', 'Massachusetts', 'Massachusetts', 'ماساتشوستس', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MI', 'Michigan', 'Michigan', 'ميشيغان', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MN', 'Minnesota', 'Minnesota', 'مينيسوتا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MS', 'Mississippi', 'Mississippi', 'ميسيسيبي', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MO', 'Missouri', 'Missouri', 'ميسوري', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'MT', 'Montana', 'Montana', 'مونتانا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NE', 'Nebraska', 'Nebraska', 'نبراسكا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NV', 'Nevada', 'Nevada', 'نيفادا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NH', 'New Hampshire', 'New Hampshire', 'نيو هامبشير', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NJ', 'New Jersey', 'New Jersey', 'نيو جيرسي', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NM', 'New Mexico', 'Nouveau-Mexique', 'نيو مكسيكو', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NY', 'New York', 'New York', 'نيويورك', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'NC', 'North Carolina', 'Caroline du Nord', 'كارولينا الشمالية', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'ND', 'North Dakota', 'Dakota du Nord', 'داكوتا الشمالية', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'OH', 'Ohio', 'Ohio', 'أوهايو', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'OK', 'Oklahoma', 'Oklahoma', 'أوكلاهوما', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'OR', 'Oregon', 'Oregon', 'أوريغون', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'PA', 'Pennsylvania', 'Pennsylvanie', 'بنسلفانيا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'RI', 'Rhode Island', 'Rhode Island', 'رود آيلاند', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'SC', 'South Carolina', 'Caroline du Sud', 'كارولينا الجنوبية', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'SD', 'South Dakota', 'Dakota du Sud', 'داكوتا الجنوبية', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'TN', 'Tennessee', 'Tennessee', 'تينيسي', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'TX', 'Texas', 'Texas', 'تكساس', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'UT', 'Utah', 'Utah', 'يوتا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'VT', 'Vermont', 'Vermont', 'فيرمونت', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'VA', 'Virginia', 'Virginie', 'فيرجينيا', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'WA', 'Washington', 'Washington', 'واشنطن', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'WV', 'West Virginia', 'Virginie-Occidentale', 'فيرجينيا الغربية', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'WI', 'Wisconsin', 'Wisconsin', 'ويسكونسن', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'WY', 'Wyoming', 'Wyoming', 'وايومنغ', 'state'),
				((SELECT id FROM countries WHERE code = 'US'), 'DC', 'District of Columbia', 'District de Columbia', 'مقاطعة كولومبيا', 'district')
				-- Insert only, no conflicts expected
			`,
			Down: `
				DELETE FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US');
			`,
		},
		{
			Version:        4,
			Description:    "Load Moroccan regions",
			RequiredSchema: 1,
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"states", "morocco"},
			Up: `
				-- Moroccan Administrative Regions with full locale support
				INSERT INTO states (country_id, code, name_en, name_fr, name_ar, type) VALUES
				((SELECT id FROM countries WHERE code = 'MA'), '01', 'Tanger-Tetouan-Al Hoceima', 'Tanger-Tétouan-Al Hoceïma', 'طنجة تطوان الحسيمة', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '02', 'Oriental', 'L''Oriental', 'الشرق', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '03', 'Fez-Meknes', 'Fès-Meknès', 'فاس مكناس', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '04', 'Rabat-Sale-Kenitra', 'Rabat-Salé-Kénitra', 'الرباط سلا القنيطرة', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '05', 'Beni Mellal-Khenifra', 'Béni Mellal-Khénifra', 'بني ملال خنيفرة', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '06', 'Casablanca-Settat', 'Casablanca-Settat', 'الدار البيضاء سطات', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '07', 'Marrakech-Safi', 'Marrakech-Safi', 'مراكش آسفي', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '08', 'Draa-Tafilalet', 'Drâa-Tafilalet', 'درعة تافيلالت', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '09', 'Souss-Massa', 'Souss-Massa', 'سوس ماسة', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '10', 'Guelmim-Oued Noun', 'Guelmim-Oued Noun', 'كلميم واد نون', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '11', 'Laayoune-Sakia El Hamra', 'Laâyoune-Sakia El Hamra', 'العيون الساقية الحمراء', 'region'),
				((SELECT id FROM countries WHERE code = 'MA'), '12', 'Dakhla-Oued Ed-Dahab', 'Dakhla-Oued Ed-Dahab', 'الداخلة وادي الذهب', 'region')
				-- Insert only, no conflicts expected
			`,
			Down: `
				DELETE FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA');
			`,
		},
		{
			Version:        5,
			Description:    "Load French regions",
			RequiredSchema: 1,
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"states", "france"},
			Up: `
				-- French Administrative Regions with full locale support
				INSERT INTO states (country_id, code, name_en, name_fr, name_ar, type) VALUES
				((SELECT id FROM countries WHERE code = 'FR'), 'ARA', 'Auvergne-Rhone-Alpes', 'Auvergne-Rhône-Alpes', 'أوفيرني رون ألب', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'BFC', 'Burgundy-Franche-Comte', 'Bourgogne-Franche-Comté', 'بورغونيا فرانش كونتيه', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'BRE', 'Brittany', 'Bretagne', 'بريتاني', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'CVL', 'Centre-Loire Valley', 'Centre-Val de Loire', 'سنتر فال دو لوار', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'COR', 'Corsica', 'Corse', 'كورسيكا', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'GES', 'Grand East', 'Grand Est', 'الشرق الكبير', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'HDF', 'Hauts-de-France', 'Hauts-de-France', 'أوت دو فرانس', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'IDF', 'Ile-de-France', 'Île-de-France', 'إيل دو فرانس', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'NOR', 'Normandy', 'Normandie', 'نورماندي', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'NAQ', 'Nouvelle-Aquitaine', 'Nouvelle-Aquitaine', 'نوفيل أكيتين', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'OCC', 'Occitania', 'Occitanie', 'أوكسيتانيا', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'PDL', 'Pays de la Loire', 'Pays de la Loire', 'بييه دو لا لوار', 'region'),
				((SELECT id FROM countries WHERE code = 'FR'), 'PAC', 'Provence-Alpes-Cote d''Azur', 'Provence-Alpes-Côte d''Azur', 'بروفانس ألب كوت دازور', 'region')
				-- Insert only, no conflicts expected
			`,
			Down: `
				DELETE FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'FR');
			`,
		},
		{
			Version:        6,
			Description:    "Load major Canadian cities",
			RequiredSchema: 1,
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"cities", "canada", "major"},
			Up: `
				-- Major Canadian Cities with comprehensive data
				INSERT INTO cities (country_id, state_id, code, name_en, name_fr, name_ar, latitude, longitude, population, timezone, is_capital) VALUES
				-- Ontario
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'toronto', 'Toronto', 'Toronto', 'تورونتو', 43.6532, -79.3832, 2930000, 'America/Toronto', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'ottawa', 'Ottawa', 'Ottawa', 'أوتاوا', 45.4215, -75.6972, 994000, 'America/Toronto', TRUE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'hamilton', 'Hamilton', 'Hamilton', 'هاملتون', 43.2557, -79.8711, 537000, 'America/Toronto', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'london', 'London', 'London', 'لندن', 42.9849, -81.2453, 383000, 'America/Toronto', FALSE),
				-- Quebec  
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'QC'), 'montreal', 'Montreal', 'Montréal', 'مونتريال', 45.5017, -73.5673, 1780000, 'America/Toronto', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'QC'), 'quebec', 'Quebec City', 'Ville de Québec', 'مدينة كيبيك', 46.8139, -71.2080, 542000, 'America/Toronto', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'QC'), 'laval', 'Laval', 'Laval', 'لافال', 45.6066, -73.7124, 422000, 'America/Toronto', FALSE),
				-- British Columbia
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'BC'), 'vancouver', 'Vancouver', 'Vancouver', 'فانكوفر', 49.2827, -123.1207, 675000, 'America/Vancouver', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'BC'), 'victoria', 'Victoria', 'Victoria', 'فيكتوريا', 48.4284, -123.3656, 91000, 'America/Vancouver', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'BC'), 'burnaby', 'Burnaby', 'Burnaby', 'بيرنابي', 49.2488, -122.9805, 232000, 'America/Vancouver', FALSE),
				-- Alberta
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'AB'), 'calgary', 'Calgary', 'Calgary', 'كالجاري', 51.0447, -114.0719, 1306000, 'America/Edmonton', FALSE),
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'AB'), 'edmonton', 'Edmonton', 'Edmonton', 'إدمونتون', 53.5461, -113.4938, 981000, 'America/Edmonton', FALSE)
				-- Insert only, no conflicts expected
			`,
			Down: `
				DELETE FROM cities WHERE country_id = (SELECT id FROM countries WHERE code = 'CA');
			`,
		},
		{
			Version:        7, 
			Description:    "Load major US cities",
			RequiredSchema: 1,
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"cities", "usa", "major"},
			Up: `
				-- Major US Cities with comprehensive data
				INSERT INTO cities (country_id, state_id, code, name_en, name_fr, name_ar, latitude, longitude, population, timezone, is_capital) VALUES
				-- New York
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'NY'), 'nyc', 'New York City', 'New York', 'نيويورك', 40.7128, -74.0060, 8336000, 'America/New_York', FALSE),
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'NY'), 'albany', 'Albany', 'Albany', 'ألباني', 42.6526, -73.7562, 97000, 'America/New_York', FALSE),
				-- California
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'CA'), 'la', 'Los Angeles', 'Los Angeles', 'لوس أنجلوس', 34.0522, -118.2437, 3980000, 'America/Los_Angeles', FALSE),
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'CA'), 'sf', 'San Francisco', 'San Francisco', 'سان فرانسيسكو', 37.7749, -122.4194, 874000, 'America/Los_Angeles', FALSE),
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'CA'), 'sd', 'San Diego', 'San Diego', 'سان دييغو', 32.7157, -117.1611, 1420000, 'America/Los_Angeles', FALSE),
				-- Texas  
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'TX'), 'houston', 'Houston', 'Houston', 'هيوستن', 29.7604, -95.3698, 2320000, 'America/Chicago', FALSE),
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'TX'), 'dallas', 'Dallas', 'Dallas', 'دالاس', 32.7767, -96.7970, 1343000, 'America/Chicago', FALSE),
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'TX'), 'austin', 'Austin', 'Austin', 'أوستن', 30.2672, -97.7431, 965000, 'America/Chicago', FALSE),
				-- Illinois
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'IL'), 'chicago', 'Chicago', 'Chicago', 'شيكاغو', 41.8781, -87.6298, 2716000, 'America/Chicago', FALSE),
				-- Maryland (Healthcare platform specific)
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'MD'), 'baltimore', 'Baltimore', 'Baltimore', 'بالتيمور', 39.2904, -76.6122, 593000, 'America/New_York', FALSE),
				-- Florida
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'FL'), 'miami', 'Miami', 'Miami', 'ميامي', 25.7617, -80.1918, 467000, 'America/New_York', FALSE),
				-- Washington
				((SELECT id FROM countries WHERE code = 'US'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'US') AND code = 'WA'), 'seattle', 'Seattle', 'Seattle', 'سياتل', 47.6062, -122.3321, 753000, 'America/Los_Angeles', FALSE)
				-- Insert only, no conflicts expected
			`,
			Down: `
				DELETE FROM cities WHERE country_id = (SELECT id FROM countries WHERE code = 'US');
			`,
		},
		{
			Version:        8,
			Description:    "Load comprehensive nationality data linked to countries",
			RequiredSchema: 5, // Requires nationalities table from schema version 5
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"nationalities", "core"},
			Up: `
				-- Simple nationality data - one nationality per country
				INSERT INTO nationalities (country_id, name_en, name_fr, name_ar, code, is_primary) VALUES
				((SELECT id FROM countries WHERE code = 'CA'), 'Canadian', 'Canadien', 'كندي', 'CA', TRUE),
				((SELECT id FROM countries WHERE code = 'US'), 'American', 'Américain', 'أمريكي', 'US', TRUE),
				((SELECT id FROM countries WHERE code = 'MA'), 'Moroccan', 'Marocain', 'مغربي', 'MA', TRUE),
				((SELECT id FROM countries WHERE code = 'FR'), 'French', 'Français', 'فرنسي', 'FR', TRUE)
				
				ON CONFLICT (country_id, code) DO UPDATE SET
					name_en = EXCLUDED.name_en,
					name_fr = EXCLUDED.name_fr,
					name_ar = EXCLUDED.name_ar,
					is_primary = EXCLUDED.is_primary,
					updated_at = CURRENT_TIMESTAMP;
			`,
			Down: `
				DELETE FROM nationalities WHERE country_id IN (
					SELECT id FROM countries WHERE code IN ('CA', 'US', 'MA', 'FR')
				);
			`,
		},
		/*
		{
			Version:        9,
			Description:    "Load comprehensive Moroccan cities for all 12 regions",
			RequiredSchema: 1,
			Environment:    "disabled", // Temporarily disabled to fix ON CONFLICT issue
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"cities", "morocco", "comprehensive"},
			Up: `
				-- Comprehensive Moroccan Cities for all 12 regions
				INSERT INTO cities (country_id, state_id, name_en, name_fr, name_ar, latitude, longitude, population, timezone, is_capital) VALUES
				
				-- Region 01: Tanger-Tétouan-Al Hoceima
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Tangier', 'Tanger', 'طنجة', 35.7595, -5.8340, 947952, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Tetouan', 'Tétouan', 'تطوان', 35.5889, -5.3626, 380787, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Al Hoceima', 'Al Hoceïma', 'الحسيمة', 35.2517, -3.9317, 56716, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Chefchaouen', 'Chefchaouen', 'شفشاون', 35.1688, -5.2636, 42786, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Larache', 'Larache', 'العرائش', 35.1932, -6.1561, 125008, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Ksar el-Kebir', 'Ksar el-Kébir', 'القصر الكبير', 35.0056, -5.9018, 126617, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Asilah', 'Asilah', 'أصيلة', 35.4658, -6.0347, 31420, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Fnideq', 'Fnideq', 'الفنيدق', 35.8500, -5.3583, 34026, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '01'), 'Ouazzane', 'Ouazzane', 'وزان', 34.7833, -5.5833, 59606, 'Africa/Casablanca', FALSE),
				
				-- Region 02: L'Oriental
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Oujda', 'Oujda', 'وجدة', 34.6814, -1.9086, 494252, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Nador', 'Nador', 'الناظور', 35.1681, -2.9287, 178540, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Berkane', 'Berkane', 'بركان', 34.9183, -2.3217, 109237, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Taourirt', 'Taourirt', 'تاوريرت', 34.4092, -2.8953, 82518, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Jerada', 'Jerada', 'جرادة', 34.3117, -2.1631, 43506, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Al Aroui', 'Al Aroui', 'العروي', 35.0167, -2.9833, 28858, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Zaio', 'Zaio', 'زايو', 34.9333, -2.7333, 31340, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '02'), 'Bouarfa', 'Bouarfa', 'بوعرفة', 32.5167, -1.9667, 27677, 'Africa/Casablanca', FALSE),
				
				-- Region 03: Fès-Meknès  
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Meknes', 'Meknès', 'مكناس', 33.8935, -5.5473, 632079, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Taza', 'Taza', 'تازة', 34.2133, -4.0100, 148456, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Errachidia', 'Errachidia', 'الراشيدية', 31.9314, -4.4241, 92374, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Sefrou', 'Sefrou', 'صفرو', 33.8275, -4.8378, 79887, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Ifrane', 'Ifrane', 'إفران', 33.5228, -5.1106, 14659, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Azrou', 'Azrou', 'أزرو', 33.4333, -5.2167, 54749, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Khenifra', 'Khénifra', 'خنيفرة', 32.9353, -5.6681, 86716, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '03'), 'Midelt', 'Midelt', 'ميدلت', 32.6833, -4.7500, 55304, 'Africa/Casablanca', FALSE),
				
				-- Region 04: Rabat-Salé-Kénitra
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'Sale', 'Salé', 'سلا', 34.0531, -6.7986, 890403, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'Kenitra', 'Kénitra', 'القنيطرة', 34.2610, -6.5802, 431282, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'Temara', 'Témara', 'تمارة', 33.9287, -6.9063, 344924, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'Khemisset', 'Khemisset', 'الخميسات', 33.8244, -6.0668, 127917, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'Skhirate', 'Skhirate', 'الصخيرات', 33.8561, -7.0374, 38026, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'Sidi Kacem', 'Sidi Kacem', 'سيدي قاسم', 34.2167, -5.7000, 75969, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '04'), 'Sidi Slimane', 'Sidi Slimane', 'سيدي سليمان', 34.2667, -5.9333, 79749, 'Africa/Casablanca', FALSE),
				
				-- Region 05: Béni Mellal-Khénifra
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'Beni Mellal', 'Béni Mellal', 'بني ملال', 32.3372, -6.3498, 192676, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'Kasba Tadla', 'Kasba Tadla', 'قصبة تادلة', 32.5975, -6.2656, 44598, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'Fkih Ben Salah', 'Fkih Ben Salah', 'الفقيه بن صالح', 32.5017, -6.6914, 78840, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'Khouribga', 'Khouribga', 'خريبكة', 32.8833, -6.9000, 196196, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'Oued Zem', 'Oued Zem', 'وادي زم', 32.8667, -6.5667, 84910, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '05'), 'Azilal', 'Azilal', 'أزيلال', 31.9667, -6.5667, 48887, 'Africa/Casablanca', FALSE),
				
				-- Region 06: Casablanca-Settat  
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'Settat', 'Settat', 'سطات', 33.0008, -7.6164, 142250, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'El Jadida', 'El Jadida', 'الجديدة', 33.2316, -8.5007, 196789, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'Berrechid', 'Berrechid', 'برشيد', 33.2650, -7.5869, 136634, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'Benslimane', 'Benslimane', 'بن سليمان', 33.6133, -7.1211, 46522, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'Azemmour', 'Azemmour', 'أزمور', 33.2833, -8.3333, 42500, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '06'), 'Sidi Bennour', 'Sidi Bennour', 'سيدي بنور', 32.6500, -8.4167, 43153, 'Africa/Casablanca', FALSE),
				
				-- Region 07: Marrakech-Safi
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'Safi', 'Safi', 'آسفي', 32.2994, -9.2372, 308508, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'Essaouira', 'Essaouira', 'الصويرة', 31.5085, -9.7595, 77966, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'Kelaa des Sraghna', 'Kelâa des Sraghna', 'قلعة السراغنة', 32.0544, -7.4059, 98943, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'Youssoufia', 'Youssoufia', 'اليوسفية', 32.2465, -8.5286, 67628, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '07'), 'Ben Guerir', 'Ben Guerir', 'بن جرير', 32.2167, -7.9333, 56806, 'Africa/Casablanca', FALSE),
				
				-- Region 08: Drâa-Tafilalet  
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'Ouarzazate', 'Ouarzazate', 'ورزازات', 30.9189, -6.8939, 71067, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'Zagora', 'Zagora', 'زاكورة', 30.3273, -5.8372, 40069, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'Tinghir', 'Tinghir', 'تنغير', 31.5147, -5.5333, 43632, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'Rissani', 'Rissani', 'الريصاني', 31.2794, -4.2594, 20469, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '08'), 'Arfoud', 'Erfoud', 'أرفود', 31.3333, -4.2333, 22930, 'Africa/Casablanca', FALSE),
				
				-- Region 09: Souss-Massa
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'Agadir', 'Agadir', 'أكادير', 30.4278, -9.5981, 421844, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'Inezgane', 'Inezgane', 'إنزكان', 30.3556, -9.5361, 101243, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'Taroudant', 'Taroudant', 'تارودانت', 30.4667, -8.8667, 80149, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'Tiznit', 'Tiznit', 'تزنيت', 29.6833, -9.7333, 74699, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'Ouled Teima', 'Ouled Teïma', 'أولاد تايمة', 30.3833, -9.2167, 71915, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '09'), 'Ait Melloul', 'Aït Melloul', 'آيت ملول', 30.3333, -9.5000, 168000, 'Africa/Casablanca', FALSE),
				
				-- Region 10: Guelmim-Oued Noun
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '10'), 'Guelmim', 'Guelmim', 'كلميم', 28.9833, -10.0500, 118318, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '10'), 'Tan-Tan', 'Tan-Tan', 'طان طان', 28.4333, -11.1000, 73209, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '10'), 'Sidi Ifni', 'Sidi Ifni', 'سيدي إفني', 29.3833, -10.1833, 20051, 'Africa/Casablanca', FALSE),
				
				-- Region 11: Laâyoune-Sakia El Hamra
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '11'), 'Laayoune', 'Laâyoune', 'العيون', 27.1500, -13.2000, 271000, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '11'), 'Es Semara', 'Es-Semara', 'السمارة', 26.7333, -11.6667, 64049, 'Africa/Casablanca', FALSE),
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '11'), 'Boujdour', 'Boujdour', 'بوجدور', 26.1333, -14.5000, 46129, 'Africa/Casablanca', FALSE),
				
				-- Region 12: Dakhla-Oued Ed-Dahab
				((SELECT id FROM countries WHERE code = 'MA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = '12'), 'Dakhla', 'Dakhla', 'الداخلة', 23.7167, -15.9333, 106277, 'Africa/Casablanca', FALSE);
			`,
			Down: `
				-- Remove comprehensive Morocco cities, keeping only basic ones
				DELETE FROM cities 
				WHERE country_id = (SELECT id FROM countries WHERE code = 'MA')
				AND name_en NOT IN ('Casablanca', 'Rabat', 'Fez', 'Marrakech');
			`,
		},
		*/
		{
			Version:        10,
			Description:    "Load comprehensive Moroccan insurance system data",
			RequiredSchema: 4, // Requires insurance tables from schema migration 4
			Environment:    "all",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"insurance", "morocco"},
			Up: `
				-- Morocco Insurance Types (comprehensive coverage)
				INSERT INTO insurance_types (country_id, code, name_en, name_fr, name_ar, is_default, sort_order) VALUES
				-- Public Health Insurance
				((SELECT id FROM countries WHERE code = 'MA'), 'amo', 'AMO (Public Health Insurance)', 'AMO (Assurance Maladie Obligatoire)', 'التأمين الصحي الإجباري (أمو)', TRUE, 1),
				((SELECT id FROM countries WHERE code = 'MA'), 'ramed', 'RAMED (Medical Assistance Regime)', 'RAMED (Régime d''Assistance Médicale)', 'نظام المساعدة الطبية (راميد)', FALSE, 2),
				
				-- Private Health Insurance
				((SELECT id FROM countries WHERE code = 'MA'), 'private', 'Private Health Insurance', 'Assurance Maladie Privée', 'التأمين الصحي الخاص', FALSE, 3),
				((SELECT id FROM countries WHERE code = 'MA'), 'complementary', 'Complementary Health Insurance', 'Assurance Maladie Complémentaire', 'التأمين الصحي التكميلي', FALSE, 4),
				
				-- Professional and Occupational
				((SELECT id FROM countries WHERE code = 'MA'), 'professional', 'Professional Insurance', 'Assurance Professionnelle', 'التأمين المهني', FALSE, 5),
				((SELECT id FROM countries WHERE code = 'MA'), 'military', 'Military Health Insurance', 'Assurance Santé Militaire', 'التأمين الصحي العسكري', FALSE, 6),
				
				-- International and Travel
				((SELECT id FROM countries WHERE code = 'MA'), 'international', 'International Health Insurance', 'Assurance Santé Internationale', 'التأمين الصحي الدولي', FALSE, 7),
				((SELECT id FROM countries WHERE code = 'MA'), 'travel', 'Travel Health Insurance', 'Assurance Voyage Santé', 'تأمين السفر الصحي', FALSE, 8),
				
				-- Other
				((SELECT id FROM countries WHERE code = 'MA'), 'other', 'Other', 'Autre', 'أخرى', FALSE, 99)
				
				ON CONFLICT (country_id, code) DO UPDATE SET
					name_en = EXCLUDED.name_en,
					name_fr = EXCLUDED.name_fr,
					name_ar = EXCLUDED.name_ar,
					is_default = EXCLUDED.is_default,
					sort_order = EXCLUDED.sort_order,
					updated_at = CURRENT_TIMESTAMP;

				-- Morocco Insurance Providers by Type
				
				-- AMO (Public Health Insurance) Providers
				INSERT INTO insurance_providers (insurance_type_id, code, name_en, name_fr, name_ar, is_default, sort_order) VALUES
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'amo'), 'cnops', 'CNOPS (National Social Security Fund)', 'CNOPS (Caisse Nationale des Organismes de Prévoyance Sociale)', 'الصندوق الوطني لمنظمات الاحتياط الاجتماعي', TRUE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'amo'), 'cnss', 'CNSS (National Social Security Fund)', 'CNSS (Caisse Nationale de Sécurité Sociale)', 'الصندوق الوطني للضمان الاجتماعي', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'amo'), 'cmr', 'CMR (Moroccan Retirement Fund)', 'CMR (Caisse Marocaine des Retraites)', 'الصندوق المغربي للتقاعد', FALSE, 3),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'amo'), 'rcar', 'RCAR (Pension Fund for Government Employees)', 'RCAR (Régime Collectif d''Allocation de Retraite)', 'نظام الأجير الجماعي للتقاعد', FALSE, 4),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'amo'), 'other_amo', 'Other AMO Provider', 'Autre Prestataire AMO', 'مقدم أمو آخر', FALSE, 99),

				-- RAMED (Medical Assistance Regime) Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'ramed'), 'ramed_ministry', 'Ministry of Health RAMED', 'RAMED Ministère de la Santé', 'راميد وزارة الصحة', TRUE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'ramed'), 'ramed_regional', 'Regional RAMED Office', 'Bureau Régional RAMED', 'مكتب راميد الجهوي', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'ramed'), 'other_ramed', 'Other RAMED Provider', 'Autre Prestataire RAMED', 'مقدم راميد آخر', FALSE, 99),

				-- Private Health Insurance Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'axa_morocco', 'AXA Assurance Maroc', 'AXA Assurance Maroc', 'أكسا التأمين المغرب', TRUE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'wafa_assurance', 'Wafa Assurance', 'Wafa Assurance', 'وفا للتأمين', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'saham_assurance', 'Saham Assurance', 'Saham Assurance', 'سهام للتأمين', FALSE, 3),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'atlanta_assurance', 'Atlanta Assurance', 'Atlanta Assurance', 'أتلانتا للتأمين', FALSE, 4),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'mamda_mcma', 'MAMDA & MCMA', 'MAMDA & MCMA', 'ماموا ومكماء', FALSE, 5),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'rma_watanya', 'RMA Watanya', 'RMA Watanya', 'الوطنية للتأمين', FALSE, 6),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'cnia_saada', 'CNIA Saada', 'CNIA Saada', 'الشركة الوطنية للاستثمار الفلاحي السعادة', FALSE, 7),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'zurich_morocco', 'Zurich Assurance Maroc', 'Zurich Assurance Maroc', 'زيوريخ التأمين المغرب', FALSE, 8),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'private'), 'other_private', 'Other Private Insurance', 'Autre Assurance Privée', 'تأمين خاص آخر', FALSE, 99),

				-- Complementary Health Insurance Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'complementary'), 'complementary_cnops', 'CNOPS Complementary', 'Complémentaire CNOPS', 'تكميلي الصندوق الوطني', FALSE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'complementary'), 'complementary_private', 'Private Complementary Insurance', 'Assurance Complémentaire Privée', 'التأمين التكميلي الخاص', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'complementary'), 'mutual_insurance', 'Mutual Insurance Company', 'Mutuelle d''Assurance', 'شركة التأمين التبادلي', FALSE, 3),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'complementary'), 'other_complementary', 'Other Complementary Insurance', 'Autre Assurance Complémentaire', 'تأمين تكميلي آخر', FALSE, 99),

				-- Professional Insurance Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'professional'), 'medical_syndicate', 'Medical Professional Syndicate', 'Syndicat Professionnel Médical', 'النقابة المهنية الطبية', FALSE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'professional'), 'lawyer_insurance', 'Bar Association Insurance', 'Assurance du Barreau', 'تأمين نقابة المحامين', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'professional'), 'engineer_insurance', 'Engineers Insurance', 'Assurance des Ingénieurs', 'تأمين المهندسين', FALSE, 3),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'professional'), 'other_professional', 'Other Professional Insurance', 'Autre Assurance Professionnelle', 'تأمين مهني آخر', FALSE, 99),

				-- Military Health Insurance Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'military'), 'far_medical', 'Royal Armed Forces Medical Service', 'Service de Santé des Forces Armées Royales', 'الخدمة الطبية للقوات المسلحة الملكية', TRUE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'military'), 'gendarmerie_medical', 'Royal Gendarmerie Medical Service', 'Service Médical de la Gendarmerie Royale', 'الخدمة الطبية للدرك الملكي', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'military'), 'other_military', 'Other Military Insurance', 'Autre Assurance Militaire', 'تأمين عسكري آخر', FALSE, 99),

				-- International Health Insurance Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'international'), 'allianz_care', 'Allianz Care', 'Allianz Care', 'أليانز كير', FALSE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'international'), 'cigna_global', 'Cigna Global', 'Cigna Global', 'سيجنا العالمية', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'international'), 'bupa_global', 'Bupa Global', 'Bupa Global', 'بوبا العالمية', FALSE, 3),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'international'), 'april_international', 'April International', 'April International', 'أبريل الدولية', FALSE, 4),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'international'), 'other_international', 'Other International Insurance', 'Autre Assurance Internationale', 'تأمين دولي آخر', FALSE, 99),

				-- Travel Health Insurance Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'travel'), 'travel_guard', 'Travel Guard', 'Travel Guard', 'ترافل جارد', FALSE, 1),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'travel'), 'world_nomads', 'World Nomads', 'World Nomads', 'ورلد نومادز', FALSE, 2),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'travel'), 'allianz_travel', 'Allianz Travel', 'Allianz Travel', 'أليانز للسفر', FALSE, 3),
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'travel'), 'other_travel', 'Other Travel Insurance', 'Autre Assurance Voyage', 'تأمين سفر آخر', FALSE, 99),

				-- Other Insurance Providers
				((SELECT id FROM insurance_types WHERE country_id = (SELECT id FROM countries WHERE code = 'MA') AND code = 'other'), 'other_provider', 'Other Insurance Provider', 'Autre Prestataire d''Assurance', 'مقدم تأمين آخر', TRUE, 1)

				ON CONFLICT (insurance_type_id, code) DO UPDATE SET
					name_en = EXCLUDED.name_en,
					name_fr = EXCLUDED.name_fr,
					name_ar = EXCLUDED.name_ar,
					is_default = EXCLUDED.is_default,
					sort_order = EXCLUDED.sort_order,
					updated_at = CURRENT_TIMESTAMP;
			`,
			Down: `
				-- Remove all Moroccan insurance data
				DELETE FROM insurance_providers 
				WHERE insurance_type_id IN (
					SELECT id FROM insurance_types 
					WHERE country_id = (SELECT id FROM countries WHERE code = 'MA')
				);
				
				DELETE FROM insurance_types 
				WHERE country_id = (SELECT id FROM countries WHERE code = 'MA');
			`,
		},
	}
}

// GetTestDataMigrations returns test/dev specific data migrations
func GetTestDataMigrations() []DataMigration {
	return []DataMigration{
		{
			Version:        100,
			Description:    "Load test healthcare entities sample data",
			RequiredSchema: 1,
			Environment:    "dev",
			CanRerun:       false, // Prevent duplicate insertions
			Tags:           []string{"test", "entities", "sample"},
			Up: `
				-- Test data migration example
				-- This would be used for development/testing environments only
				INSERT INTO cities (country_id, state_id, code, name_en, name_fr, name_ar, latitude, longitude, population, timezone) VALUES
				((SELECT id FROM countries WHERE code = 'CA'), (SELECT id FROM states WHERE country_id = (SELECT id FROM countries WHERE code = 'CA') AND code = 'ON'), 'test_city', 'Test City', 'Ville Test', 'مدينة التجربة', 45.0, -75.0, 50000, 'America/Toronto');
			`,
			Down: `
				DELETE FROM cities WHERE code = 'test_city';
			`,
		},
	}
}