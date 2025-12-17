package router

import (
	"github.com/gin-gonic/gin"

	"task_manager/controllers"
	"task_manager/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// -----------------------------
	// PUBLIC ROUTES (NO AUTH REQUIRED)
	// -----------------------------
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// -----------------------------
	// AUTHENTICATED ROUTES
	// -----------------------------
	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthRequired())

	// READ-ONLY TASK ROUTES (USER + ADMIN)
	authRoutes.GET("/tasks", controllers.GetAllTasks)
	authRoutes.GET("/tasks/:id", controllers.GetTaskByID)

	// -----------------------------
	// ADMIN-ONLY ROUTES
	// -----------------------------
	adminRoutes := authRoutes.Group("/")
	adminRoutes.Use(middleware.AdminRequired())

	// Admin manages tasks
	adminRoutes.POST("/tasks", controllers.CreateTask)
	adminRoutes.PUT("/tasks/:id", controllers.UpdateTask)
	adminRoutes.DELETE("/tasks/:id", controllers.DeleteTask)

	// Promote user to admin
	adminRoutes.POST("/promote/:id", controllers.PromoteUser)

	return r
}
