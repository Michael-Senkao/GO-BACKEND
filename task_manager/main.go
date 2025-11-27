package main

import "task_manager/router"

func main() {
	r := router.InitRoutes()
	r.Run("localhost:8080")
}
