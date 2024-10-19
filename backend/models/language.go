package models

type Language struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:255;not null;uniqueIndex:idx_language_extension"`
	Extension string `gorm:"size:50;not null;uniqueIndex:idx_language_extension"`
	Supported bool   `gorm:"default:false"`
}
