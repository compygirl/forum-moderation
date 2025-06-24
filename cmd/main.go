// @title Forum API
// @version 1.0
// @description This is a sample forum API server written in Go.
// @host https://localhost:8082
// @BasePath /

package main

import (
	"context"
	"forum/cmd/config"
	"forum/internal/server"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configObj := config.CreateConfig() // new config was created with port 8082
	if err := config.ReadConfig("cmd/config/Config.json", configObj); err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background() // You can also use context.WithCancel, context.WithTimeout, etc., as needed
	srv := server.InitServer(configObj, ctx)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
