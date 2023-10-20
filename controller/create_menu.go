package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateMenu(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Memeriksa apakah pengguna memiliki izin untuk membuat menu (isAdmin=true)
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Hanya Admin yang dapat menambahkan"})
		}

		// Menguraikan data menu dari JSON yang diterima
		var Menu struct {
			NameRestoran string  `json:"name_restoran"`
			Name         string  `json:"name"`
			Location     string  `json:"location"`
			Description  string  `json:"description"`
			Price        float32 `json:"price"`
			Category     string  `json:"category"`
			AvalibleMenu int     `json:"available_menu"`
		}

		if err := c.Bind(&Menu); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Membuat Menu baru dan mengaitkannya dengan pengguna
		newMenu := model.Menu{
			NameRestoran: Menu.NameRestoran,
			Name:         Menu.Name,
			Location:     Menu.Location,
			Description:  Menu.Description,
			Price:        Menu.Price,
			Category:     Menu.Category,
			AvalibleMenu: Menu.AvalibleMenu,
			UserID:       user.ID, // Mengaitkan menu dengan pengguna yang membuatnya
		}

		if err := db.Create(&newMenu).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create menu"})
		}

		// Mengembalikan respons sukses jika berhasil
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":    false,
			"message":  "menu created successfully",
			"menuData": newMenu, // Mengirim data menu yang baru saja dibuat
		})
	}
}
