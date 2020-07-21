package application

import (
	"go-bestflight/application/web/http"
	"go-bestflight/domain/services/routeservice"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/database"
	"go-bestflight/resources/file"
	"log"
)

func Start(filePath string, port string) {
	database.Connect()
	cache.Connect()
	file.Sync(filePath)

	routesFromFile, err := file.ReadFile()
	if err != nil {
		log.Fatalf("could not read from file %s: %v", filePath, err)
	}

	routeservice.LoadRoutes(routesFromFile)
	http.Start(port)
	// cli.StartAdvisor()
}
