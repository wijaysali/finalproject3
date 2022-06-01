package models

type SocialMedia struct {
	ID uint `gorm:"PrimaryKey" json:"id"`
	DefaultModel
	Name           string `gorm:"not null"`
	SocialMediaUrl string `gorm:"not null"`
	UserID         uint   `json:"user_id"`
}
