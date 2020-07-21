package http

import (
	"context"
	"fmt"
	"go-bestflight/application/web/http/routes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	server *http.Server
)

// Start the http server.
func Start(port string, mode string, loggerWriter io.Writer) {
	gin.SetMode(mode)
	router := gin.New()

	if loggerWriter == nil {
		router.Use(gin.Logger())
	} else {
		router.Use(gin.LoggerWithWriter(loggerWriter))
	}

	routes.InscribeRoutes(router)

	server = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start the HTTP server: %v", err)
		}
	}()

	message := fmt.Sprintf("HTTP server running on port %s...", port)
	fmt.Println(message)
	log.Println(message)
}

// GracefullShutdown allows gracefull shutdown.
func GracefullShutdown(quitChan chan os.Signal) {
	go func() {
		log.Println("gracefull shutdown enabled")

		oscall := <-quitChan

		log.Printf("system call:%+v", oscall)
		log.Println("shutting server down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("erro when shuttingdown the server:", err)
		}

		log.Println("finished")
		log.Println("press ctrl + D to exit program completely")
		fmt.Println("\nserver shutted down. Press ctrl + D to exit program completely")
	}()
}
