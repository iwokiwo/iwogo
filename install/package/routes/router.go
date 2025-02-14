package routes

import (
	"iwogo/auth"
	"iwogo/helper"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	//auth service
	authService := auth.NewService()

	//setting gin
	router := gin.Default()

	router.Use(helper.SetTimeZone(os.Getenv("asia/makassar")))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"POST, GET, OPTIONS, PUT, DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "X-CSRF-Token"}
	router.Use(cors.New(config))
	// public routes
	router.Static("/storage", "./storage")
	//router.Use(static.Serve("/storage", static.LocalFile("/storage", false)))
	api := router.Group("/api/v1")

	// new routing
	// UserRouter(db, api, authService)
	// UnitRouter(db, api, authService)
	// RfidRouter(db, api, authService)
	// List of route functions

	routeFunctions := []func(*gorm.DB, *gin.RouterGroup, auth.Service) *gin.RouterGroup{
		UserRouter,
	}

	// Loop through and call each route function
	for _, routeFunc := range routeFunctions {
		_ = routeFunc(db, api, authService) // Call the function and discard return value
	}

	return router
}
