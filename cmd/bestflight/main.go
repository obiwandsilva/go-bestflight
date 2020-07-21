package main

import (
	"go-bestflight/application"
	"os"
)

func main() {
	filePath := os.Args[1]
	port := os.Args[2]

	application.Start(filePath, port)
}
