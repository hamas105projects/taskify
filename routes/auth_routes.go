package routes

import (
	"taskify/usecase"

	"github.com/gin-gonic/gin"
)

// AuthRoutes mengatur rute-rute yang berkaitan dengan autentikasi
func AuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", usecase.Register)
		auth.POST("/login", usecase.Login)
	}
}
