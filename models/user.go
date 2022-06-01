package models

type User struct {
	ID uint `gorm:"PrimaryKey" json:"id"`
	DefaultModel
	Username    string `gorm:"not null;unique" validate:"required"`
	Email       string `gorm:"not null;unique" validate:"required,email"`
	Password    string `gorm:"not null" validate:"required,min=6"`
	Age         uint   `gorm:"not null" validate:"required,gt=8"`
	Photos      []Photo
	SocialMedia []SocialMedia
	Comments    []Comment
}
