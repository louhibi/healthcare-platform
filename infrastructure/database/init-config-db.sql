-- Config Service Database Initialization
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

INSERT INTO config_settings (key, value, is_public, description) VALUES
    ('app.name', 'Healthcare Platform', true, 'Application display name'),
    ('app.version', '0.1.0', true, 'Application version')
ON CONFLICT (key) DO NOTHING;

INSERT INTO feature_flags (name, enabled, description, is_public) VALUES
    ('appointments.enabled', true, 'Enable appointment booking', true),
    ('registration.enabled', true, 'Allow public self-registration', true)
ON CONFLICT (name) DO NOTHING;
