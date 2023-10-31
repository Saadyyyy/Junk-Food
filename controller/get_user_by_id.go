package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetUserById(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Mendapatkan ID pengguna yang diminta dari parameter URL
		requestedUserID := c.Param("id")

		// Mengambil ID pengguna dari database berdasarkan username yang terkait dengan token
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Membandingkan apakah ID pengguna yang diminta sesuai dengan ID pengguna yang ditemukan dalam token
		if requestedUserID != fmt.Sprint(user.ID) {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Access denied"})
		}

		// Mengembalikan data pengguna dalam format yang diinginkan
		userData := map[string]interface{}{
			"id":           user.UserID,
			"username":     user.Username,
			"phone_number": user.PhoneNumber,
			"email":        user.Email,
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"error": false, "message": "User data retrieved successfully", "user": userData})
	}
}
