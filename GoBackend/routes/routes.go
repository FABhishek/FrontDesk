package routes

import (
	"frontdesk/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine,
	queriesHandler handlers.QueriesHandler) {
	v1 := router.Group("/api/v1")
	{
		queries := v1.Group("/queries")
		{
			queries.POST("", queriesHandler.CreateQuery)
			queries.GET("", queriesHandler.GetQueries)
			queries.PATCH("/:id/resolve", queriesHandler.ResolveQuery)
			queries.GET("/faqs", queriesHandler.GetFAQs)
		}
	}
}
