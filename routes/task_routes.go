package routes

import (
	"taskify/middlewares" // Pastikan import middlewares
	"taskify/usecase"

	"github.com/gin-gonic/gin"
)

// TaskRoutes mengatur rute-rute yang berkaitan dengan tugas
func TaskRoutes(api *gin.RouterGroup) {
	// Group untuk endpoint yang memerlukan otentikasi
	authenticated := api.Group("/")
	authenticated.Use(middlewares.AuthMiddleware()) // Terapkan AuthMiddleware

	{
		// Rute Tasks di bawah Project
		authenticated.POST("/projects/:project_id/tasks", usecase.CreateTask)
		authenticated.GET("/projects/:project_id/tasks", usecase.GetTasksByProject)
		authenticated.GET("/projects/:project_id/tasks/:task_id", usecase.GetTaskByID)
		authenticated.PUT("/projects/:project_id/tasks/:task_id", usecase.UpdateTask)
		authenticated.DELETE("/projects/:project_id/tasks/:task_id", usecase.DeleteTask)
	}
}
