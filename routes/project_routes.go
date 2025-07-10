package routes

import (
	"taskify/middlewares" // Pastikan import middlewares
	"taskify/usecase"

	"github.com/gin-gonic/gin"
)

// ProjectRoutes mengatur rute-rute yang berkaitan dengan proyek
func ProjectRoutes(api *gin.RouterGroup) {
	// Group untuk endpoint yang memerlukan otentikasi
	authenticated := api.Group("/")
	authenticated.Use(middlewares.AuthMiddleware()) // Terapkan AuthMiddleware

	{
		authenticated.POST("/projects", usecase.CreateProject)
		authenticated.GET("/projects", usecase.GetProjects)
		authenticated.GET("/projects/detail/:id", usecase.GetProjectByID)
		authenticated.PUT("/projects/detail/:id", usecase.UpdateProject)
		authenticated.DELETE("/projects/detail/:id", usecase.DeleteProject)
	}
}
