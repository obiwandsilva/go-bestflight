package routes

import (
	controllers "go-bestflight/application/web/http/controllers/routecontroller"

	"github.com/gin-gonic/gin"
)

// InscribeRoutes ...
func InscribeRoutes(server *gin.Engine) {
	server.POST("/routes", controllers.AddNewRoute)
	server.GET("/routes", controllers.BestRoute)
}
