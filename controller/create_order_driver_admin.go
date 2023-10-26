package controller

import (
	"Junk-Food/middleware"
	"Junk-Food/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateOrderDriverByAdmin(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the token from the Authorization header
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Check if the Authorization header contains a Bearer token
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token format. Use 'Bearer [token]'"})
		}

		// Extract the token from the header
		tokenString = tokenString[7:]

		// Verify the token
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Get the user information associated with the token
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Check if the user has admin privileges (IsAdmin=true)
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Hanya Admin yang dapat menambahkan"})
		}

		// Parse JSON data from the request body
		var order struct {
			Message          string `json:"message"`
			PickUpLocation   string `json:"pick-up_location"`
			DeliveryLocation string `json:"delivery_location"`
			DetailOrderIDs   []uint `json:"detail_order_ids"` // Add DetailOrderIDs field to the struct
		}

		if err := c.Bind(&order); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Create a new DriverOrderAdmin instance
		newOrder := model.DriverOrderAdmin{
			Message:          order.Message,
			PickUpLocation:   order.PickUpLocation,
			DeliveryLocation: order.DeliveryLocation,
		}

		// Link the DetailOrder records by their IDs
		for _, detailOrderID := range order.DetailOrderIDs {
			var detailOrder model.DetailOrder
			if err := db.First(&detailOrder, detailOrderID).Error; err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid DetailOrder ID"})
			}
			newOrder.DetailOrder = append(newOrder.DetailOrder, detailOrder)
		}

		if err := db.Create(&newOrder).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create order"})
		}

		// Return a success response if successful
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":     false,
			"status":    "Order created successfully",
			"orderData": newOrder, // Send the newly created order data
		})
	}
}
