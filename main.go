// Package main is the main package for the application
package main

import (
	"github.com/brave/go-sync/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	server.StartServer()
}
