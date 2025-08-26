package main

import (
    "database/sql"
    "errors"
)

type ConfigService struct { db *sql.DB }

func NewConfigService(db *sql.DB) *ConfigService { return &ConfigService{db: db} }

func (s *ConfigService) GetPublicSettings() (map[string]string, error) {
    rows, err := s.db.Query("SELECT key, value FROM config_settings WHERE is_public = true")
    if err != nil { return nil, err }
    defer rows.Close()
    res := map[string]string{}
    for rows.Next() { var k,v string; if err := rows.Scan(&k,&v); err!=nil { return nil, err }; res[k]=v }
    return res, nil
}

func (s *ConfigService) GetPublicFlags() (map[string]bool, error) {
    rows, err := s.db.Query("SELECT name, enabled FROM feature_flags WHERE is_public = true")
    if err != nil { return nil, err }
    defer rows.Close()
    res := map[string]bool{}
    for rows.Next() { var n string; var e bool; if err := rows.Scan(&n,&e); err!=nil { return nil, err }; res[n]=e }
    return res, nil
}

func (s *ConfigService) GetSetting(key string) (*ConfigSetting, error) {
    row := s.db.QueryRow("SELECT id, key, value, is_public, description, updated_at FROM config_settings WHERE key=$1", key)
    var cs ConfigSetting; var desc sql.NullString
    if err := row.Scan(&cs.ID, &cs.Key, &cs.Value, &cs.IsPublic, &desc, &cs.UpdatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) { return nil, nil }
        return nil, err
    }
    if desc.Valid { cs.Description = &desc.String }
    return &cs, nil
}

func (s *ConfigService) UpsertSetting(key, value string, isPublic bool, description *string) (*ConfigSetting, error) {
    row := s.db.QueryRow(`INSERT INTO config_settings (key, value, is_public, description, updated_at) VALUES ($1,$2,$3,$4,CURRENT_TIMESTAMP)
        ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value, is_public=EXCLUDED.is_public, description=EXCLUDED.description, updated_at=CURRENT_TIMESTAMP
        RETURNING id, key, value, is_public, description, updated_at`, key, value, isPublic, description)
    var cs ConfigSetting; var desc sql.NullString
    if err := row.Scan(&cs.ID, &cs.Key, &cs.Value, &cs.IsPublic, &desc, &cs.UpdatedAt); err != nil { return nil, err }
    if desc.Valid { cs.Description=&desc.String }
    return &cs, nil
}

func (s *ConfigService) ListSettings(includePrivate bool) ([]ConfigSetting, error) {
    query := "SELECT id, key, value, is_public, description, updated_at FROM config_settings"
    if !includePrivate { query += " WHERE is_public = true" }
    rows, err := s.db.Query(query)
    if err != nil { return nil, err }
    defer rows.Close()
    var list []ConfigSetting
    for rows.Next() { var cs ConfigSetting; var desc sql.NullString; if err := rows.Scan(&cs.ID,&cs.Key,&cs.Value,&cs.IsPublic,&desc,&cs.UpdatedAt); err!=nil { return nil, err }; if desc.Valid { cs.Description=&desc.String }; list=append(list, cs) }
    return list, nil
}

func (s *ConfigService) ListFeatureFlags(includePrivate bool) ([]FeatureFlag, error) {
    query := "SELECT id, name, enabled, description, is_public, updated_at FROM feature_flags"
    if !includePrivate { query += " WHERE is_public = true" }
    rows, err := s.db.Query(query)
    if err != nil { return nil, err }
    defer rows.Close()
    var list []FeatureFlag
    for rows.Next() { var ff FeatureFlag; var desc sql.NullString; if err := rows.Scan(&ff.ID,&ff.Name,&ff.Enabled,&desc,&ff.IsPublic,&ff.UpdatedAt); err!=nil { return nil, err }; if desc.Valid { ff.Description=&desc.String }; list=append(list, ff) }
    return list, nil
}

func (s *ConfigService) UpsertFeatureFlag(name string, enabled bool, isPublic bool, description *string) (*FeatureFlag, error) {
    row := s.db.QueryRow(`INSERT INTO feature_flags (name, enabled, is_public, description, updated_at) VALUES ($1,$2,$3,$4,CURRENT_TIMESTAMP)
        ON CONFLICT (name) DO UPDATE SET enabled=EXCLUDED.enabled, is_public=EXCLUDED.is_public, description=EXCLUDED.description, updated_at=CURRENT_TIMESTAMP
        RETURNING id, name, enabled, description, is_public, updated_at`, name, enabled, isPublic, description)
    var ff FeatureFlag; var desc sql.NullString
    if err := row.Scan(&ff.ID,&ff.Name,&ff.Enabled,&desc,&ff.IsPublic,&ff.UpdatedAt); err!=nil { return nil, err }
    if desc.Valid { ff.Description=&desc.String }
    return &ff, nil
}
