package main

import (
	"calendar/internal/handler"
	"calendar/internal/service"
	"flag"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	r := gin.New()

	r.Use(handler.RequestLogger())

	r.StaticFile("/", "./index.html")

	eventService := service.NewService()
	eventHandler := handler.NewHandler(eventService)

	eventHandler.RegisterRoutes(r)

	r.Run(":" + *port)
}
