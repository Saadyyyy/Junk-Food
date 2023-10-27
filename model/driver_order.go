package model

import (
	"gorm.io/gorm"
)

type DriverOrderAdmin struct {
	gorm.Model
	Message          string
	PickUpLocation   string        `json:"pick-up_location"`
	DeliveryLocation string        `json:"delivery_location"`
	DetailOrder      []DetailOrder `gorm:"foreignKey:DriverOrderAdminID" json:"detail_order"`
}

type DriverOrderUser struct {
	gorm.Model
	Message          string `json:"message"`
	PickUpLocation   string `json:"pick-up_location"`
	DeliveryLocation string `json:"delivery_location"`
	UserID           uint   `json:"user_id"`
}
