package main

import (
	"app/cmd/app"
	_ "app/docs"
)

// @title App Service
// @version 1.0
// @description application service

// @host localhost:8081
func main() {
	app.Run()
}
