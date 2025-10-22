package main

import (
	"frontdesk/config"
	db "frontdesk/database"
	"frontdesk/handlers"
	"frontdesk/repositories"
	"frontdesk/routes"
	"frontdesk/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default() // initialize the router for gin.
	config.LoadConfig()     // load the configurations.
	db.Initialize()

	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow the frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Allow cookies to be sent with cross-origin requests
	}))

	//dependency injection
	queriesRepository := repositories.NewQueriesRepository(db.DB)
	queriesService := services.NewQueriesService(queriesRepository)
	queriesHandler := handlers.NewQueriesHandler(queriesService)

	// route setup
	routes.SetupRoutes(
		router,
		queriesHandler,
	)

	if err := router.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
