package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func EditUser(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Mendapatkan ID pengguna dari parameter URL
		userID := c.Param("id")

		// Mencari pengguna berdasarkan ID
		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "User not found"})
		}

		// Memeriksa apakah pengguna yang diautentikasi memiliki hak akses atau kode otorisasi yang sesuai
		if user.Username != username {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Unauthorized to edit this user"})
		}

		// Mengekstrak data yang ingin diubah dari body permintaan PUT
		var updateUser model.User
		if err := c.Bind(&updateUser); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Mengubah data yang sesuai dalam entitas pengguna
		if updateUser.Username != "" {
			user.Username = updateUser.Username
		}
		if updateUser.Email != "" {
			user.Email = updateUser.Email
		}
		if updateUser.PhoneNumber != "" {
			user.PhoneNumber = updateUser.PhoneNumber
		}
		if updateUser.IsAdmin != user.IsAdmin {
			user.IsAdmin = updateUser.IsAdmin
		}

		// Menyimpan perubahan ke basis data
		if err := db.Save(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to update user"})
		}

		// Menghapus password dari respons
		user.Password = ""

		// Mengirim respons berhasil
		return c.JSON(http.StatusOK, map[string]interface{}{"error": false, "message": "User updated successfully", "user": user})
	}
}
