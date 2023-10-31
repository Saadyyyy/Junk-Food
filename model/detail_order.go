package model

import (
	"time"
)

type DetailOrder struct {
	DetailOrderID      uint
	MenuID             uint       `json:"event_id"` // ID acara yang makanannya terkait
	UserID             uint       `json:"user_id"`  // ID pengguna yang membeli makanan
	Quantity           int        `json:"quantity"` // Jumlah makanan yang dibeli
	KodeVoucher        string     `json:"kode_voucher"`
	TotalCost          int        `json:"total_cost"`
	InvoiceNumber      string     `json:"invoice_number"` // Nomor invoice untuk makanan
	DriverOrderAdminID uint       // Add a foreign key to establish the relationship
	CreatedAt          *time.Time `json:"created_at"` // Kolom created_at yang diharapkan tipe data *time.Time
	UpdatedAt          time.Time
}

type DetailOrderAdmin struct {
	DetailOrderAdminID uint
	DriverOrderAdminID uint   `json:"admin_id"`       // ID pengguna yang membeli makanan
	InvoiceNumber      string `json:"invoice_number"` // Nomor invoice untuk makanan
	// DriverOrderAdminID uint   // Add a foreign key to establish the relationship
	CreatedAt *time.Time `json:"created_at"` // Kolom created_at yang diharapkan tipe data *time.Time
	UpdatedAt time.Time
}
type DetailOrderUser struct {
	DetailOrderUserID uint
	DriverOrderUserID uint   `json:""`               // ID pengguna yang membeli makanan
	InvoiceNumber     string `json:"invoice_number"` // Nomor invoice untuk makanan
	// DriverOrderAdminID uint   // Add a foreign key to establish the
	CreatedAt *time.Time `json:"created_at"` // Kolom created_at yang diharapkan tipe data *time.Time
	UpdatedAt time.Time
}
