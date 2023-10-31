package model

import (
	"time"
)

type DriverOrderAdmin struct {
	DriverOrderAdminID uint
	Message            string
	PickUpLocation     string        `json:"pick-up_location"`
	DeliveryLocation   string        `json:"delivery_location"`
	DetailOrder        []DetailOrder `gorm:"foreignKey:DriverOrderAdminID" json:"detail_order"`
	CreatedAt          *time.Time    `json:"created_at"` // Kolom created_at yang diharapkan tipe data *time.Time
	UpdatedAt          time.Time
}

type DriverOrderUser struct {
	DriverOrderUserID uint
	Message           string     `json:"message"`
	PickUpLocation    string     `json:"pick-up_location"`
	DeliveryLocation  string     `json:"delivery_location"`
	UserID            uint       `json:"user_id"`
	CreatedAt         *time.Time `json:"created_at"` // Kolom created_at yang diharapkan tipe data *time.Time
	UpdatedAt         time.Time
}
