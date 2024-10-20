package models

type Repo struct {
	ID              uint   `gorm:"primaryKey"`
	Owner           string `gorm:"size:255;not null;uniqueIndex:idx_owner_name"`
	Name            string `gorm:"size:255;not null;uniqueIndex:idx_owner_name"`
	ScanComplete    bool   `gorm:"default:false"`
	InstallStarted  bool   `gorm:"default:false"`
	InstallComplete bool   `gorm:"default:false"`
	RescanSecurity  bool   `gorm:"default:false"`
}

type RepoLanguage struct {
	ID           uint     `gorm:"primaryKey;autoIncrement"`
	RepositoryID uint     `gorm:"not null;index"`
	LanguageID   uint     `gorm:"not null;index"`
	Language     Language `gorm:"foreignKey:LanguageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Repository   Repo     `gorm:"foreignKey:RepositoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
