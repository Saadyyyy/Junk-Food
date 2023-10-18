package main

import (
	"Junk-Food/config"
	"log"
)

func main() {
	// Setup router
	router := config.SetupRouter()

	// Mulai server Echo pada alamat dan port tertentu (misalnya, :8080)
	err := router.Start(":8001")
	if err != nil {
		log.Fatal(err)
	}
}
