package config

import (
	"Junk-Food/model"
	"Junk-Food/routes"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

func InitializeDatabase(config DatabaseConfig) (*gorm.DB, error) {
	// Konfigurasi koneksi database MySQL dengan GORM
	dsn := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DBName + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate tabel pengguna
	db.AutoMigrate(&model.User{})
	// db.AutoMigrate(&model.Event{})
	// db.AutoMigrate(&model.Order{})
	// db.AutoMigrate(&model.OrderItem{})
	// db.AutoMigrate(&model.Ticket{})
	// db.AutoMigrate(&model.Promo{})

	return db, nil
}

func SetupRouter() *echo.Echo {
	// Inisialisasi database
	dbConfig := DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		Username: "saady",
		Password: "853211",
		DBName:   "junk_food",
	}

	db, err := InitializeDatabase(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Membuat instance router Echo
	router := echo.New()

	// Middleware untuk log
	router.Use(middleware.Logger())

	// Menggunakan routes yang telah dipisahkan
	routes.SetupRoutes(router, db)

	return router
}
