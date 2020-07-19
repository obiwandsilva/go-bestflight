package http

import (
	"fmt"
	"go-bestflight/application/web/http/routes"

	"github.com/gin-gonic/gin"
)

// Start starts http server.
func Start(port string) {
	server := gin.Default()

	routes.InscribeRoutes(server)
	server.Run(fmt.Sprintf(":%s", port))
}
