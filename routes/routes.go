package routes

import (
	"Junk-Food/controller"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// e.Use(Logger())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := []byte(os.Getenv("SECRET_JWT"))

	// Menggunakan routes yang telah dipisahkan
	e.POST("/signup", controller.Signup(db, secretKey))
	e.POST("/signin", controller.Signin(db, secretKey))
	e.GET("/user/:id", controller.GetUserById(db, secretKey))
	e.PUT("/user/change-password/:id", controller.ChangePassword(db, secretKey))
	e.POST("/menu/create", controller.CreateMenu(db, secretKey))
	e.GET("/menu", controller.GetMenus(db, secretKey))
	e.GET("/menu/:id", controller.GetMenuByID(db, secretKey))
	// e.GET("/user/:id", controllers.GetUserProfile(db, secretKey))
	// e.POST("/event/create", controllers.CreateEvent(db, secretKey))
	// e.GET("/event", controllers.GetEvents(db, secretKey))
	// e.GET("/event/:id", controllers.GetEventByID(db, secretKey))
	// e.PUT("/user/:id", controllers.EditUser(db, secretKey))
	// e.PUT("/admin/user/:id", controllers.EditUserByAdmin(db, secretKey))
	e.DELETE("/admin/user/:id", controller.DeleteMenuByAdmin(db, secretKey))
	// e.GET("/admin/user", controllers.GetAllUsersByAdmin(db, secretKey))
	e.POST("/user/buy", controller.BuyMenu(db, secretKey))
	// e.GET("/user/ticket", controllers.GetTicketsByUser(db, secretKey))
	e.POST("/admin/promo", controller.CreatePromo(db, secretKey))
	// e.GET("/user/promo", controllers.GetPromos(db, secretKey))
	e.GET("/user/detailOrder/:invoiceNumber", controller.GetDetailTrascation(db, secretKey))
	makan := controller.NewMakanUsecase()
	e.POST("/chatbot/recommend-makan", func(c echo.Context) error {
		return controller.RecommendMakan(c, makan)
	})
	e.POST("/admin/list-order-admin", controller.CreateOrderDriverByAdmin(db, secretKey))
	e.POST("/user/list-order-user", controller.CreateOrderDriverByUser(db, secretKey))
	e.GET("/driver/list-admin", controller.GetListOrderByAdmin(db, secretKey))
	e.GET("/driver/list-user", controller.GetListOrderByUser(db, secretKey))

}
