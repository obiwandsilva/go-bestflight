package http

import (
	"context"
	"fmt"
	"go-bestflight/application/web/http/routes"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	server *http.Server
)

// Start the http server.
func Start(port string) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		router := gin.Default()
		routes.InscribeRoutes(router)

		server = &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: router,
		}

		wg.Done()

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("could not start the HTTP server: %v", err)
		}

		fmt.Println(fmt.Sprintf("HTTP server running on port %s...", port))
	}()

	wg.Wait()
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	log.Println("Shuting server down...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server failed to shutdown:", err)
	}
}
