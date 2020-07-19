package routes

import (
	"go-bestflight/application/web/http/controllers"

	"github.com/gin-gonic/gin"
)

// InscribeRoutes ...
func InscribeRoutes(server *gin.Engine) {
	server.GET("/routes/best", controllers.BestRoute)
}
