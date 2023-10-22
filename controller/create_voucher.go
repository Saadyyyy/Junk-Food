package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreatePromo(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Memeriksa apakah pengguna memiliki izin untuk membuat promo (isAdmin=true)
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Hanya Admin yang dapat menambahkan promo"})
		}

		// Menguraikan data promo dari JSON yang diterima
		var promo struct {
			Name                 string `json:"name"`
			KodeVoucher          string `json:"kode_voucher"`
			JumlahPotonganPersen int    `json:"jumlah_potongan_persen"`
		}

		if err := c.Bind(&promo); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Membuat promo baru
		newPromo := model.Voucher{
			Name:                 promo.Name,
			KodeVoucher:          promo.KodeVoucher,
			JumlahPotonganPersen: promo.JumlahPotonganPersen,
		}

		if err := db.Create(&newPromo).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create promo"})
		}

		// Mengembalikan respons sukses jika berhasil
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":     false,
			"message":   "Promo created successfully",
			"promoData": newPromo, // Mengirim data promo yang baru saja dibuat
		})
	}
}
