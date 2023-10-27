package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ChangePassword(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Middleware Autentikasi
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memverifikasi token
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mendapatkan ID pengguna dari parameter URL
		userID := c.Param("id")

		// Mengambil ID pengguna dari database berdasarkan username yang terkait dengan token
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "User not found"})
		}

		// Memeriksa apakah ID pengguna dalam token sesuai dengan ID pengguna dalam parameter URL
		if userID != fmt.Sprint(user.ID) {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Access denied"})
		}

		// Mendapatkan data permintaan perubahan kata sandi dari body
		var req model.ChangePasswordRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Validasi data permintaan (currentPassword, newPassword, dll.)

		// Verifikasi kata sandi saat ini
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Incorrect current password"})
		}

		// Hash dan simpan kata sandi baru
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to hash new password"})
		}

		// Update kata sandi pengguna di database
		user.Password = string(hashedNewPassword)
		db.Save(&user)

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Password updated successfully"})
	}
}
