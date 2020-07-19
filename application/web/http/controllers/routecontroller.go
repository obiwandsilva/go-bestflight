package controllers

import (
	"go-bestflight/services/routeservice"

	"github.com/gin-gonic/gin"
)

// BestRoute ...
func BestRoute(ctx *gin.Context) {
	boarding := ctx.Query("board")
	destination := ctx.Query("dest")
	bestRoute := routeservice.BestRoute(boarding, destination)

	ctx.JSON(200, gin.H{
		"result": bestRoute,
	})
}
