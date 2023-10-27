package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	IsAdmin     bool   `gorm:"default:false" json:"is_admin"`  // Menambahkan nilai default
	IsDriver    bool   `gorm:"default:false" json:"is_driver"` // Menambahkan nilai default

	Menu []Menu `gorm:"foreignKey:UserID" json:"Menu"`
}

// Buat struct untuk permintaan perubahan kata sandi
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}
