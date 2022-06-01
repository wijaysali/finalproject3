package models

type Comment struct {
	ID uint `gorm:"PrimaryKey" json:"id"`
	DefaultModel
	UserID  uint   `json:"user_id"`
	PhotoID uint   `json:"photo_id" validate:"required"`
	Message string `gorm:"not null;" validate:"required" json:"message"`
}
