package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetPromos(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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
			// Token tidak valid tetapkan username menjadi string kosong
			username = ""
		}

		// Mengambil daftar voucher dari database
		var vouchers []model.Voucher
		if err := db.Find(&vouchers).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch vouchers"})
		}

		// Mengembalikan daftar voucher dalam format yang diinginkan
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":    false,
			"username": username, // Mengirimkan username pengguna (kosong jika token tidak valid)
			"vouchers": vouchers,
		})
	}
}
