package models

type Photo struct {
	ID uint `gorm:"PrimaryKey" json:"id"`
	DefaultModel
	Title    string `gorm:"not null;" validate:"required" json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `gorm:"not null;" validate:"required" json:"photo_url"`
	UserID   uint   `json:"user_id"`
	Comments []Comment
}
