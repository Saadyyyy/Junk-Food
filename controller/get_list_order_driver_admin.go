package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetListOrderByAdmin(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Mengambil informasi user yang terkait dengan token
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Memeriksa apakah admin yang diautentikasi memiliki status IsDriver yang true
		if !user.IsDriver {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Anda bukan admin, tidak bisa menghapus data user lain!"})
		}

		// Mengambil daftar list dari database dengan detail orders
		var orders []model.DriverOrderAdmin
		if err := db.Preload("DetailOrder").Find(&orders).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch orders with details"})
		}

		// Mengembalikan daftar list dalam format yang diinginkan
		return c.JSON(http.StatusOK, map[string]interface{}{"error": false, "List_Orders": orders})
	}
}
