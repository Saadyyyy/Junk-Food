package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func BuyMenu(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Menguraikan data pembelian tiket dari JSON yang diterima
		var menuPurchase struct {
			MenuID      uint   `json:"menu_id"`
			Quantity    int    `json:"quantity"`
			KodeVoucher string `json:"kode_voucher"`
		}

		if err := c.Bind(&menuPurchase); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Memeriksa apakah menu yang akan dibeli tiketnya ada
		var menu model.Menu
		menuResult := db.First(&menu, menuPurchase.MenuID)
		if menuResult.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "menu not found"})
		}

		// Mencari promo berdasarkan kode voucher yang dimasukkan (jika ada)
		var promo model.Voucher
		if menuPurchase.KodeVoucher != "" {
			promoResult := db.Where("kode_voucher = ?", menuPurchase.KodeVoucher).First(&promo)
			if promoResult.Error != nil {
				// Jika promo tidak ditemukan, kirimkan pesan kesalahan
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid voucher code"})
			}
		}

		// Menghitung total biaya pembelian tiket
		totalCost := int(menu.Price) * menuPurchase.Quantity

		// Jika promo ditemukan, menghitung potongan biaya tiket
		if promo.ID != 0 {
			potongan := (float64(promo.JumlahPotonganPersen) / 100) * float64(totalCost)
			totalCost -= int(potongan)
		}

		// Membuat entri baru dalam tabel menu
		detailOrder := model.DetailOrder{
			MenuID:        menu.ID,
			UserID:        user.ID,
			Quantity:      menuPurchase.Quantity,
			TotalCost:     totalCost,                // Total biaya tiket setelah potongan
			InvoiceNumber: generateInvoiceNumber(),  // Simpan nomor invoice dalam tiket
			KodeVoucher:   menuPurchase.KodeVoucher, // Menyimpan kode voucher dalam entri tiket
		}

		if err := db.Create(&detailOrder).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create menu"})
		}

		// Mengurangi jumlah tiket yang tersedia
		menu.DecrementMenus(menuPurchase.Quantity)
		db.Save(&menu)

		// Menyiapkan respons JSON yang mencakup kode voucher yang digunakan
		responseData := map[string]interface{}{
			"error":          false,
			"message":        "menu purchased successfully",
			"detailOrderID":  detailOrder.ID,            // Mengirimkan ID tiket yang telah dibeli
			"invoice_number": detailOrder.InvoiceNumber, // Mengirimkan nomor invoice
			"totalCost":      totalCost,                 // Mengirimkan total biaya tiket
			"kode_voucher":   menuPurchase.KodeVoucher,  // Mengirimkan kode voucher yang digunakan
		}

		return c.JSON(http.StatusOK, responseData)
	}
}

func generateInvoiceNumber() string {
	// Menggunakan waktu sekarang dan nomor acak untuk membuat nomor invoice unik
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(1000) // Ganti dengan rentang nomor yang sesuai
	invoiceNumber := fmt.Sprintf("%d-%d", timestamp, randomNum)
	return invoiceNumber
}
