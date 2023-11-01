package routes

import (
	"Junk-Food/controller"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// e.Use(Logger())
	godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	secretKey := []byte(os.Getenv("SECRET_JWT"))

	// Menggunakan routes yang telah dipisahkan
	e.POST("/signup", controller.Signup(db, secretKey))
	e.POST("/signin", controller.Signin(db, secretKey))
	e.PUT("/user/change-password/:id", controller.ChangePassword(db, secretKey))
	e.POST("/menu/create", controller.CreateMenu(db, secretKey))
	e.GET("/menu", controller.GetMenus(db, secretKey))
	e.GET("/menu/:id", controller.GetMenuByID(db, secretKey))
	e.PUT("/user/:id", controller.EditUser(db, secretKey))
	e.DELETE("/admin/menu/:id", controller.DeleteMenuByAdmin(db, secretKey))
	e.POST("/user/buy", controller.BuyMenu(db, secretKey))
	e.POST("/admin/voucher", controller.CreateVoucher(db, secretKey))
	e.GET("/user/voucher", controller.GetVouchers(db, secretKey))
	e.GET("/user/detailOrder/:invoiceNumber", controller.GetDetailTrascation(db, secretKey))
	e.GET("/user/detailOrder/id", controller.GetDetailTrascationByID(db, secretKey))
	makan := controller.NewMakanUsecase()
	e.POST("/chatbot/recommend-makan", func(c echo.Context) error {
		return controller.RecommendMakan(c, makan)
	})
	e.POST("/admin/list-order-admin", controller.CreateOrderDriverByAdmin(db, secretKey))
	e.POST("/user/list-order-user", controller.CreateOrderDriverByUser(db, secretKey))
	e.GET("/driver/list-admin", controller.GetListOrderByAdmin(db, secretKey))
	e.GET("/driver/list-user", controller.GetListOrderByUser(db, secretKey))
	e.DELETE("/admin/voucher/:id", controller.DeleteVoucherByAdmin(db, secretKey))
	e.POST("/driver/admin", controller.TakeOrderAdmin(db, secretKey))
	e.POST("/driver/user", controller.TakeOrderUser(db, secretKey))

}
