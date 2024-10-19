// models/user.go
package models

// User represents the user model
type User struct {
	ID             uint   `gorm:"primaryKey"`
	GithubUsername string `gorm:"unique;not null"`
	HasAccess      bool   `gorm:"default:false"`
}
