package models

import (
	"time"
)

type Base struct {
	Id        uint       `gorm:"unique;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
