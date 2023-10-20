package model

import "time"

type DetailOrder struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	MenuID        uint       `json:"event_id"` // ID acara yang makanannya terkait
	UserID        uint       `json:"user_id"`  // ID pengguna yang membeli makanan
	Quantity      int        `json:"quantity"` // Jumlah makanan yang dibeli
	KodeVoucher   string     `json:"kode_voucher"`
	TotalCost     int        `json:"total_cost"`
	InvoiceNumber string     `json:"invoice_number"` // Nomor invoice untuk makanan
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     time.Time
}
