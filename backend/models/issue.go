package models

type SecurityIssue struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	RepositoryID    uint   `gorm:"not null;index"`
	Path            string `gorm:"size:255;not null"`
	StartLine       int    `gorm:"not null"`
	EndLine         int    `gorm:"not null"`
	GithubNumber    int    `gorm:"not null"`
	FullDescription string `gorm:"not null"`
	FixSuggested    bool   `gorm:"not nulldefault:false"`
	Repository      Repo   `gorm:"foreignKey:RepositoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
