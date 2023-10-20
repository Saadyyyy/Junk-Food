package model

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	NameRestoran string  `json:"name_restoran"`
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	Category     string  `json:"category"`
	AvalibleMenu int     `json:"available_menu"`
	UserID       uint    `json:"user_id"` // ID pengguna yang membuat menu
}

// Metode untuk mengurangi jumlah Makanan yang tersedia
func (e *Menu) DecrementMenus(quantity int) {
	if e.AvalibleMenu >= quantity {
		e.AvalibleMenu -= quantity
	}
}
