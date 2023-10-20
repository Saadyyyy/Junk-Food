package model

import "time"

type Order struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `json:"user_id"`    // ID pengguna yang melakukan pesanan
	TotalCost  int        `json:"total_cost"` // Total biaya pesanan
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  time.Time
	OrderItems []OrderItem // Detail pesanan yang berisi makanan yang dibeli
}

type OrderItem struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	OrderID  uint `json:"order_id"`   // ID pesanan yang makanannya terkait
	Voucher  uint `json:"voucher_id"` // ID makanan yang dibeli
	Quantity int  `json:"quantity"`   // Jumlah makanan yang dibeli dalam pesanan ini
	Subtotal int  `json:"subtotal"`   // Total biaya untuk makanan ini dalam pesanan ini
}
