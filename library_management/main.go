package main

import (
	"fmt"

	"library_management/controllers"
	"library_management/services"
)

func main() {
	lib := services.NewLibrary()
	// optional: seed sample data
	lib.SeedSampleData()

	ctrl := controllers.NewController(lib)

	fmt.Println("Welcome to the Library Management System")
	ctrl.Start()
}
