package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateOrderDriverByUser(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Menguraikan data order dari JSON yang diterima
		var order struct {
			Message          string `json:"message"`
			PickUpLocation   string `json:"pick-up_location"`
			DeliveryLocation string `json:"delivery_location"`
			UserID           uint   `json:"user_id"`
		}

		if err := c.Bind(&order); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Membuat order baru dan mengaitkannya dengan pengguna
		newOrder := model.DriverOrderUser{
			Message:          order.Message,
			PickUpLocation:   order.PickUpLocation,
			DeliveryLocation: order.DeliveryLocation,
			UserID:           user.UserID,
		}

		if err := db.Create(&newOrder).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create menu"})
		}

		// Mengembalikan respons sukses jika berhasil
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":     false,
			"status":    "order created successfully",
			"orderData": newOrder, // Mengirim data menu yang baru saja dibuat
		})
	}
}
