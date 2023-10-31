package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func TakeOrderAdmin(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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
			AdminID uint `json:"list_order"`
		}

		if err := c.Bind(&menuPurchase); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Memeriksa apakah menu yang akan dibeli tiketnya ada
		var menu model.DriverOrderAdmin
		menuResult := db.First(&menu, menuPurchase.AdminID)
		if menuResult.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "menu not found"})
		}
		// Membuat entri baru dalam tabel menu
		detailOrder := model.DetailOrderAdmin{
			DriverOrderAdminID: menuPurchase.AdminID,    // Total biaya tiket setelah potongan
			InvoiceNumber:      generateInvoiceNumber(), // Simpan nomor invoice dalam tiket
		}
		if err := db.Create(&detailOrder).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create menu"})
		}

		// Mengurangi jumlah tiket yang tersedia
		// menu.DecrementMenus(menuPurchase.Quantity)
		// db.Save(&menu)

		// Menyiapkan respons JSON yang mencakup kode voucher yang digunakan
		responseData := map[string]interface{}{
			"message":        "Order purchased successfully",
			"detailOrderID":  detailOrder.DetailOrderAdminID, // Mengirimkan ID tiket yang telah dibeli
			"invoice_number": detailOrder.InvoiceNumber,      // Mengirimkan nomor invoice
		}

		return c.JSON(http.StatusOK, responseData)
	}
}

func TakeOrderUser(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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
			AdminID uint `json:"list_order"`
		}

		if err := c.Bind(&menuPurchase); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Memeriksa apakah menu yang akan dibeli tiketnya ada
		var menu model.DriverOrderUser
		menuResult := db.First(&menu, menuPurchase.AdminID)
		if menuResult.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "menu not found"})
		}
		// Membuat entri baru dalam tabel menu
		detailOrder := model.DetailOrderAdmin{
			DriverOrderAdminID: menuPurchase.AdminID,    // Total biaya tiket setelah potongan
			InvoiceNumber:      generateInvoiceNumber(), // Simpan nomor invoice dalam tiket
		}
		if err := db.Create(&detailOrder).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create menu"})
		}

		// Mengurangi jumlah tiket yang tersedia
		// menu.DecrementMenus(menuPurchase.Quantity)
		// db.Save(&menu)

		// Menyiapkan respons JSON yang mencakup kode voucher yang digunakan
		responseData := map[string]interface{}{
			"message":        "Order purchased successfully",
			"detailOrderID":  detailOrder.DetailOrderAdminID, // Mengirimkan ID tiket yang telah dibeli
			"invoice_number": detailOrder.InvoiceNumber,      // Mengirimkan nomor invoice
		}

		return c.JSON(http.StatusOK, responseData)
	}
}
