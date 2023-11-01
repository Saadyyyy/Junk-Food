package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetDetailTrascationByID(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memeriksa apakah header Authorization mengandung token Bearer
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token format. Use 'Bearer [token]'"})
		}

		// Ekstrak token dari header
		tokenString = tokenString[7:]

		// Memverifikasi token
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mendapatkan informasi pengguna dari basis data
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Memeriksa apakah pengguna memiliki status admin
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Access forbidden for non-admin users"})
		}

		// Mendapatkan nomor invoice dari parameter URL
		invoiceNumber := c.Param("id")

		// Mencari tiket berdasarkan nomor invoice
		var DetailOrder model.DetailOrder
		result = db.Where("user_id = ?", invoiceNumber).First(&DetailOrder)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "DetailOrder not found"})
		}

		// Mengambil detail menu berdasarkan MenuID yang ada pada detail
		var menu model.Menu
		menuResult := db.First(&menu, DetailOrder.MenuID)
		if menuResult.Error != nil {
			// Handle jika menu tidak ditemukan
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch menu data"})
		}

		// Membuat respons dengan detail pembelian tiket
		detailOrders := map[string]interface{}{
			"detail_order_ID": DetailOrder.ID,
			"user_id":         DetailOrder.UserID,
			"menu_id":         DetailOrder.MenuID,
			"menu_name":       menu.Name,
			"quantity":        DetailOrder.Quantity,
			"total_cost":      DetailOrder.TotalCost,
			"invoice_number":  DetailOrder.InvoiceNumber,
			"kode_voucher":    DetailOrder.KodeVoucher, // Menambahkan kode voucher ke respons jika ada
		}

		// Mengembalikan respons dengan detail pembelian tiket
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":            false,
			"message":          "DetailOrder details retrieved successfully",
			"DetailOrder_data": detailOrders,
		})
	}
}
