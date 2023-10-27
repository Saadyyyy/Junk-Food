package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DeleteVoucherByAdmin(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Middleware Autentikasi
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

		// Memverifikasi token dan mendapatkan informasi admin yang diautentikasi
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mendapatkan data admin dari token
		var adminUser model.User
		result := db.Where("username = ?", username).First(&adminUser)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "Anda bukan admin !"})
		}

		// Memeriksa apakah admin yang diautentikasi memiliki status IsAdmin yang true
		if !adminUser.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Anda bukan admin, tidak bisa menghapus data user lain!"})
		}

		// Mendapatkan ID pengguna dari parameter URL
		voucherID := c.Param("id")

		// Mencari pengguna berdasarkan ID
		var voucher model.Voucher
		result = db.First(&voucher, voucherID)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "voucher not found"})
		}

		// Menghapus pengguna dari basis data
		db.Delete(&voucher)

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "voucher deleted successfully"})
	}
}
