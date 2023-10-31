package model

import "time"

type Menu struct {
	MenuID       uint
	NameRestoran string     `json:"name_restoran"`
	Name         string     `json:"name"`
	Location     string     `json:"location"`
	Description  string     `json:"description"`
	Price        float32    `json:"price"`
	Category     string     `json:"category"`
	AvalibleMenu int        `json:"available_menu"`
	UserID       uint       `json:"user_id"`    // ID pengguna yang membuat menu
	CreatedAt    *time.Time `json:"created_at"` // Kolom created_at yang diharapkan tipe data *time.Time
	UpdatedAt    time.Time
}

// Metode untuk mengurangi jumlah Makanan yang tersedia
func (e *Menu) DecrementMenus(quantity int) {
	if e.AvalibleMenu >= quantity {
		e.AvalibleMenu -= quantity
	}
}
