package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`    // ID pengguna yang melakukan pesanan
	TotalCost  int         `json:"total_cost"` // Total biaya pesanan
	OrderItems []OrderItem // Detail pesanan yang berisi makanan yang dibeli
}

type OrderItem struct {
	gorm.Model
	OrderID   uint `json:"order_id"`   // ID pesanan yang makanannya terkait
	VoucherID uint `json:"voucher_id"` // ID makanan yang dibeli
	Quantity  int  `json:"quantity"`   // Jumlah makanan yang dibeli dalam pesanan ini
	Subtotal  int  `json:"subtotal"`   // Total biaya untuk makanan ini dalam pesanan ini
}
