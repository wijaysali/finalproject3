package models

import "time"

type DefaultModel struct {
	CreatedAt time.Time `gorm:"type:date" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:date" json:"updated_at"`
}
