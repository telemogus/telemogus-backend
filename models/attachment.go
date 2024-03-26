package models

type Attachment struct {
	Base
	MessageID uint `gorm:"foreignKey:MessageID"`
	FileType  string
	FilePath  string
	FileSize  int64
}
