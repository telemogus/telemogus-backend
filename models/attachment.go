package models

import (
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	MessageID uint `gorm:"foreignKey:MessageID"`
	FileType  string
	FilePath  string
	FileSize  int64
}
