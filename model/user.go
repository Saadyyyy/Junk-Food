package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	IsAdmin     bool   `gorm:"default:false" json:"isAdmin"` // Menambahkan nilai default
	// Events      []Event `gorm:"foreignKey:UserID" json:"events"`
}
