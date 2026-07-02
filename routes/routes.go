package routes

import (
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/signup", handlers.ShowSignup)

	router.POST("/signup", handlers.Signup)

	router.GET("/login", handlers.ShowLogin)

	router.POST("/login", handlers.Login)

	router.GET("/logout", handlers.Logout)

	router.GET("/dashboard", handlers.Dashboard)

}
