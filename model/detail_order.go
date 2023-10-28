package model

import (
	"gorm.io/gorm"
)

type DetailOrder struct {
	gorm.Model
	MenuID             uint   `json:"event_id"` // ID acara yang makanannya terkait
	UserID             uint   `json:"user_id"`  // ID pengguna yang membeli makanan
	Quantity           int    `json:"quantity"` // Jumlah makanan yang dibeli
	KodeVoucher        string `json:"kode_voucher"`
	TotalCost          int    `json:"total_cost"`
	InvoiceNumber      string `json:"invoice_number"` // Nomor invoice untuk makanan
	DriverOrderAdminID uint   // Add a foreign key to establish the relationship
}

type DetailOrderAdmin struct {
	gorm.Model
	DriverOrderAdminID uint   `json:"admin_id"`       // ID pengguna yang membeli makanan
	InvoiceNumber      string `json:"invoice_number"` // Nomor invoice untuk makanan
	// DriverOrderAdminID uint   // Add a foreign key to establish the relationship
}
type DetailOrderUser struct {
	gorm.Model
	DriverOrderUserID uint   `json:""`               // ID pengguna yang membeli makanan
	InvoiceNumber     string `json:"invoice_number"` // Nomor invoice untuk makanan
	// DriverOrderAdminID uint   // Add a foreign key to establish the relationship
}
