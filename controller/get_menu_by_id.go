package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetMenuByID(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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
		_, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mengambil menu ID dari path parameter
		menuIDParam := c.Param("id")

		// Mengonversi menuIDParam ke tipe data uint
		menuID, err := strconv.ParseUint(menuIDParam, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid menu ID"})
		}

		// Mencari menu berdasarkan menu ID
		var menu model.Menu
		if err := db.First(&menu, menuID).Error; err != nil {
			// menu tidak ditemukan
			if err == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "menu not found"})
			}
			// Terjadi kesalahan lain saat mengambil data dari database
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch menu"})
		}

		// Mengembalikan menu yang ditemukan dalam format yang diinginkan
		return c.JSON(http.StatusOK, map[string]interface{}{"error": false, "menu": menu})
	}
}
