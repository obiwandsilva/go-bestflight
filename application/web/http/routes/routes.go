package routes

import (
	"go-bestflight/application/web/http/controllers"

	"github.com/gin-gonic/gin"
)

// InscribeRoutes ...
func InscribeRoutes(server *gin.Engine) {
	server.POST("/routes", controllers.AddNewRoute)
	server.GET("/routes", controllers.BestRoute)
}
