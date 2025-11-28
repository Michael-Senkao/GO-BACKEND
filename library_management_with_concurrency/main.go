package main

import (
	"fmt"

	"library_management/controllers"
	"library_management/services"
)

func main() {
	lib := services.NewLibrary()
	lib.SeedSampleData()

	ctrl := controllers.NewController(lib)

	fmt.Println("Welcome to the Library Management System (Concurrent Reservations Demo)")
	fmt.Println("Note: Reservations auto-cancel after 5 seconds if not borrowed.")
	fmt.Println()
	ctrl.Start()
}
