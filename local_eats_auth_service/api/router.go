package api

import (
	_ "auth_serice/api/docs"
	"auth_serice/api/handlers"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterApi @title Auth
// @version 1.0
// @description Auth service
// @host localhost:8087
// @BasePath /
// @in header
func RouterApi(userStorage *handlers.Handler) *gin.Engine {
	router := gin.Default()
	h := handlers.NewHandler(userStorage)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	auth := router.Group("/api/auth_service")
	{
		auth.POST("/register", h.RegisterHandler)
		auth.POST("/login", h.LoginHandler)

	}

	return router
}
