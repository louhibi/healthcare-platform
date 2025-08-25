-- Sample Healthcare Entities for Multi-Tenant Platform
-- This script creates sample hospitals, clinics, and doctor offices across Canada, USA, Morocco, and France

-- Insert Healthcare Entities
INSERT INTO healthcare_entities (name, type, country, address, city, state, postal_code, phone, email, website, license, tax_id, timezone, language, currency) VALUES
-- Canada
('Toronto General Hospital', 'hospital', 'Canada', '200 Elizabeth Street', 'Toronto', 'ON', 'M5G 2C4', '(416) 340-4800', 'info@tgh.ca', 'https://www.uhn.ca/TorontoGeneral', 'ON-HOSP-001', '12345678901', 'America/Toronto', 'en', 'CAD'),
('Montreal Heart Institute', 'hospital', 'Canada', '5000 Rue Bélanger', 'Montreal', 'QC', 'H1T 1C8', '(514) 376-3330', 'info@icm-mhi.org', 'https://www.icm-mhi.org', 'QC-HOSP-002', '12345678902', 'America/Toronto', 'fr', 'CAD'),
('Vancouver Family Clinic', 'clinic', 'Canada', '1234 Robson Street', 'Vancouver', 'BC', 'V6E 1N5', '(604) 555-0123', 'contact@vanfamilyclinic.ca', 'https://www.vanfamilyclinic.ca', 'BC-CLI-001', '12345678903', 'America/Vancouver', 'en', 'CAD'),
('Dr. Sarah Johnson Family Practice', 'doctor_office', 'Canada', '789 Queen Street West', 'Toronto', 'ON', 'M6J 1G1', '(416) 555-0456', 'office@drjohnson.ca', 'https://www.drjohnson.ca', 'ON-DOC-001', '12345678904', 'America/Toronto', 'en', 'CAD'),

-- USA
('Johns Hopkins Hospital', 'hospital', 'USA', '1800 Orleans Street', 'Baltimore', 'MD', '21287', '(410) 955-5000', 'info@jhmi.edu', 'https://www.hopkinsmedicine.org', 'MD-HOSP-001', '52-0595110', 'America/New_York', 'en', 'USD'),
('Mayo Clinic', 'hospital', 'USA', '200 First Street SW', 'Rochester', 'MN', '55905', '(507) 284-2511', 'info@mayo.edu', 'https://www.mayoclinic.org', 'MN-HOSP-002', '41-0693611', 'America/Chicago', 'en', 'USD'),
('Silicon Valley Medical Group', 'clinic', 'USA', '2490 Hospital Drive', 'Mountain View', 'CA', '94040', '(650) 555-0789', 'contact@svmg.com', 'https://www.svmg.com', 'CA-CLI-001', '94-1234567', 'America/Los_Angeles', 'en', 'USD'),
('Dr. Michael Rodriguez Cardiology', 'doctor_office', 'USA', '1500 SW 1st Avenue', 'Miami', 'FL', '33129', '(305) 555-0321', 'office@drrodriguez.com', 'https://www.drrodriguez.com', 'FL-DOC-001', '65-0987654', 'America/New_York', 'en', 'USD'),

-- Morocco
('Centre Hospitalier Universitaire Ibn Rochd', 'hospital', 'Morocco', 'Avenue des Hôpitaux', 'Casablanca', 'Casablanca-Settat', '20360', '+212 522-48-80-80', 'contact@chu-casablanca.ma', 'https://www.chu-casablanca.ma', 'MA-HOSP-001', '3348123', 'Africa/Casablanca', 'ar', 'MAD'),
('Hôpital Universitaire International Cheikh Khalifa', 'hospital', 'Morocco', 'Had Soualem', 'Casablanca', 'Casablanca-Settat', '27182', '+212 522-77-77-77', 'info@huickk.ma', 'https://www.huickk.ma', 'MA-HOSP-002', '3348124', 'Africa/Casablanca', 'fr', 'MAD'),
('Clinique Atlas', 'clinic', 'Morocco', '2 Rue Oukaimeden', 'Rabat', 'Rabat-Salé-Kénitra', '10000', '+212 537-77-77-77', 'contact@cliniqueatlas.ma', 'https://www.cliniqueatlas.ma', 'MA-CLI-001', '3348125', 'Africa/Casablanca', 'fr', 'MAD'),
('Cabinet Dr. Youssef Bennani', 'doctor_office', 'Morocco', '45 Boulevard Mohammed V', 'Marrakech', 'Marrakech-Safi', '40000', '+212 524-44-44-44', 'cabinet@drbennani.ma', 'https://www.drbennani.ma', 'MA-DOC-001', '3348126', 'Africa/Casablanca', 'ar', 'MAD'),

-- France
('Hôpital Pitié-Salpêtrière', 'hospital', 'France', '47-83 Boulevard de l''Hôpital', 'Paris', 'Île-de-France', '75013', '+33 1 42 16 00 00', 'contact@psl.aphp.fr', 'https://www.aphp.fr', 'FR-HOSP-001', '77567227300548', 'Europe/Paris', 'fr', 'EUR'),
('Centre Hospitalier Universitaire de Lyon', 'hospital', 'France', '103 Grande Rue de la Croix-Rousse', 'Lyon', 'Auvergne-Rhône-Alpes', '69317', '+33 4 72 07 17 17', 'contact@chu-lyon.fr', 'https://www.chu-lyon.fr', 'FR-HOSP-002', '26690451900507', 'Europe/Paris', 'fr', 'EUR'),
('Clinique Mutualiste', 'clinic', 'France', '4 Rue Eric Tabarly', 'Nantes', 'Pays de la Loire', '44277', '+33 2 51 72 33 33', 'accueil@clinique-mutualiste.fr', 'https://www.clinique-mutualiste.fr', 'FR-CLI-001', '44242674400019', 'Europe/Paris', 'fr', 'EUR'),
('Cabinet Dr. Marie Dubois', 'doctor_office', 'France', '15 Avenue Victor Hugo', 'Nice', 'Provence-Alpes-Côte d''Azur', '06000', '+33 4 93 87 87 87', 'secretariat@drdubois.fr', 'https://www.drdubois.fr', 'FR-DOC-001', '32950394200025', 'Europe/Paris', 'fr', 'EUR');

-- Insert Sample Users for each entity
-- Note: All passwords are hashed for 'admin123'
-- The healthcare_entity_id will be replaced with actual IDs after healthcare_entities are created

-- Toronto General Hospital Users
INSERT INTO users (email, password_hash, first_name, last_name, role, healthcare_entity_id, license_number, specialization, preferred_language) VALUES
('admin.tgh@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'John', 'Administrator', 'admin', 1, NULL, NULL, 'en'),
('dr.smith.tgh@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. James', 'Smith', 'doctor', 1, 'ON-12345', 'Cardiology', 'en'),
('nurse.wilson.tgh@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Emily', 'Wilson', 'nurse', 1, 'RN-67890', 'ICU', 'en'),

-- Montreal Heart Institute Users
('admin.mhi@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Marie', 'Administrateur', 'admin', 2, NULL, NULL, 'fr'),
('dr.leblanc.mhi@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Pierre', 'Leblanc', 'doctor', 2, 'QC-54321', 'Cardiologie', 'fr'),
('infirmiere.martin.mhi@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Sophie', 'Martin', 'nurse', 2, 'OIIQ-12345', 'Soins Intensifs', 'fr'),

-- Johns Hopkins Hospital Users
('admin.jh@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Robert', 'Administrator', 'admin', 5, NULL, NULL, 'en'),
('dr.johnson.jh@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Sarah', 'Johnson', 'doctor', 5, 'MD-98765', 'Neurology', 'en'),
('nurse.brown.jh@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Michael', 'Brown', 'nurse', 5, 'RN-11111', 'Emergency', 'en'),

-- CHU Casablanca Users
('admin.chu@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Ahmed', 'Administrateur', 'admin', 9, NULL, NULL, 'ar'),
('dr.alami.chu@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Fatima', 'Alami', 'doctor', 9, 'MA-67890', 'طب الأطفال', 'ar'),
('infirmier.hassan.chu@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Hassan', 'Benali', 'nurse', 9, 'MA-INF-001', 'العناية المركزة', 'ar'),

-- Hôpital Pitié-Salpêtrière Users
('admin.psl@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Pierre', 'Administrateur', 'admin', 13, NULL, NULL, 'fr'),
('dr.martin.psl@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Claire', 'Martin', 'doctor', 13, 'FR-54321', 'Neurologie', 'fr'),
('infirmiere.dubois.psl@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Isabelle', 'Dubois', 'nurse', 13, 'IDE-67890', 'Réanimation', 'fr');

-- Add more sample users for other entities...
INSERT INTO users (email, password_hash, first_name, last_name, role, healthcare_entity_id, license_number, specialization, preferred_language) VALUES
-- Vancouver Family Clinic
('admin.vfc@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Lisa', 'Chen', 'admin', 3, NULL, NULL, 'en'),
('dr.patel.vfc@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Raj', 'Patel', 'doctor', 3, 'BC-11111', 'Family Medicine', 'en'),

-- Mayo Clinic
('dr.anderson.mayo@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Lisa', 'Anderson', 'doctor', 6, 'MN-22222', 'Oncology', 'en'),
('nurse.garcia.mayo@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Maria', 'Garcia', 'nurse', 6, 'RN-33333', 'Oncology', 'en'),

-- Clinique Atlas Morocco
('admin.atlas@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Youssef', 'Alaoui', 'admin', 11, NULL, NULL, 'fr'),
('dr.benali.atlas@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Aicha', 'Benali', 'doctor', 11, 'MA-44444', 'Médecine Générale', 'fr'),

-- CHU Lyon
('admin.lyon@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Jean', 'Dupont', 'admin', 14, NULL, NULL, 'fr'),
('dr.bernard.lyon@healthcare.local', '$2a$10$8U6ZZkVurkotqAm7TjIDgOWhOsp2U1N/i3YeiA3zqjWXYqvESbcA.', 'Dr. Anne', 'Bernard', 'doctor', 14, 'FR-55555', 'Chirurgie', 'fr');