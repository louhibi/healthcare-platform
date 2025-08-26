package main

import "time"

type ConfigSetting struct {
    ID int `json:"id"`
    Key string `json:"key"`
    Value string `json:"value"`
    IsPublic bool `json:"is_public"`
    Description *string `json:"description,omitempty"`
    UpdatedAt time.Time `json:"updated_at"`
}

type FeatureFlag struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Enabled bool `json:"enabled"`
    Description *string `json:"description,omitempty"`
    IsPublic bool `json:"is_public"`
    UpdatedAt time.Time `json:"updated_at"`
}

type BootstrapConfig struct {
    Environment string `json:"environment"`
}
