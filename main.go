package main

import (
	"system/middleware"
	"system/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")
	routes.PublicRoutes(router)
	router.Use(middleware.Auth())
	routes.UserRoutes(router)
	router.Use(middleware.Admin())
	routes.AdminRoutes(router)

	router.Run(":")
}
