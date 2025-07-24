package models

import (
	"time"
)

// Config represents the configuration model
type Config struct {
	Key       string    `json:"key" gorm:"primaryKey;size:128;not null"`
	Value     string    `json:"value" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConfigRequest represents the request payload for setting config
type ConfigRequest struct {
	Key   string `json:"key" binding:"required,max=128"`
	Value string `json:"value" binding:"required"`
}

// ConfigResponse represents the response structure
type ConfigResponse struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Config to ConfigResponse
func (c *Config) ToResponse() ConfigResponse {
	return ConfigResponse{
		Key:       c.Key,
		Value:     c.Value,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
