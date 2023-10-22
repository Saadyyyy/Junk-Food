package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetOrderItemsByUserID(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memverifikasi token
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mengambil informasi user yang terkait dengan token
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Mengambil order_items berdasarkan ID user
		var orderItems []model.OrderItem
		db.Preload("Order").Where("order.user_id = ?", user.ID).Find(&orderItems)

		// Menyiapkan data yang akan direspons
		responseOrderItems := make([]map[string]interface{}, len(orderItems))
		for i, orderItem := range orderItems {
			responseOrderItems[i] = map[string]interface{}{
				"id":        orderItem.ID,
				"order_id":  orderItem.OrderID,
				"vocher_id": orderItem.VoucherID,
				"quantity":  orderItem.Quantity,
				"subtotal":  orderItem.Subtotal,
			}
		}

		// Mengembalikan daftar order_items berdasarkan ID user
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":       false,
			"message":     "Order items by user",
			"user_id":     user.ID,
			"order_items": responseOrderItems,
		})
	}
}
