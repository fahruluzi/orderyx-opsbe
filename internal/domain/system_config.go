package domain

import "time"

type SystemConfig struct {
	Key         string      `json:"key" gorm:"primaryKey"`
	Value       interface{} `json:"value" gorm:"type:jsonb;serializer:json"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
