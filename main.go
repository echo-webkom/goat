package main

import (
	"github.com/echo-webkom/goat/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() // Debug
	s := server.New()
	server.ServeWithShutdown(s)
}
