package main

import (
	"go-bestflight/application"
	"os"
	"os/signal"
)

func main() {
	filePath := os.Args[1]
	port := os.Args[2]

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt)

	application.Start(filePath, port, quitChan)
}
