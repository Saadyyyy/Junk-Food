package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// create user
func Signin(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user model.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// Mengecek apakah username ada dalam database
		var existingUser model.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check username"})
			}
		}

		// Membandingkan password yang dimasukkan dengan password yang di-hash
		err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
		}

		// Generate JWT token
		tokenString, err := middleware.GenerateToken(existingUser.Username, secretKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		// Menyertakan ID pengguna dalam respons
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Login successful", "token": tokenString, "id": existingUser.ID})
	}
}

func Signup(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user model.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// Mengecek apakah username sudah ada dalam database
		var existingUser model.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check username"})
		}

		// Mengecek apakah email sudah ada dalam database
		result = db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check email"})
		}

		// Mengecek apakah phone number sudah ada dalam database
		result = db.Where("phone_number = ?", user.PhoneNumber).First(&existingUser)
		if result.Error == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Phone number already exists"})
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check phone number"})
		}

		// Meng-hash password dengan bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		// Menyimpan data pengguna ke database
		user.Password = string(hashedPassword)
		db.Create(&user)

		// Hapus password dari struct
		user.Password = ""

		// Generate JWT token
		tokenString, err := middleware.GenerateToken(user.Username, secretKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		// Menyertakan ID pengguna dalam response
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "User created successfully", "token": tokenString, "id": user.ID})
	}
}
