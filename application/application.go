package application

import (
	"go-bestflight/application/cli"
	"go-bestflight/application/web/http"
	"go-bestflight/domain/services/routeservice"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/database"
	"go-bestflight/resources/file"
	"io"
	"log"
	"os"
)

func configLogFile(filePath string) io.Writer {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Println("could not configure log file")
	}

	log.SetOutput(file)

	return file
}

func Start(filePath string, port string, quitChan chan os.Signal) {
	loggerWriter := configLogFile("info.log")
	database.Connect()
	cache.Connect()
	file.Sync(filePath)

	routesFromFile, err := file.ReadFile()
	if err != nil {
		log.Fatalf("could not read from file %s: %v", filePath, err)
	}

	routeservice.LoadRoutes(routesFromFile)
	http.Start(port, "release", loggerWriter)
	http.GracefullShutdown(quitChan)
	cli.StartAdvisor()
}
