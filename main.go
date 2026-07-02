package main

import (
	"auth/config"
	"auth/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte("super-secret-key"))

	router.Use(sessions.Sessions("auth-session", store))

	routes.SetupRoutes(router)

	router.Run(":8080")
}
