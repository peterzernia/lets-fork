package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/peterzernia/lets-fork/restaurant"
	"github.com/peterzernia/lets-fork/websocket"
)

func main() {
	hub := websocket.NewHub()

	go hub.Run()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	config.AllowOrigins = []string{"*"}

	router.Use(cors.New(config))

	api := router.Group("/api/v1")
	restaurant.InitializeRoutes(api)
	websocket.InitializeRoutes(api, hub)
	api.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	port := ":" + os.Getenv("PORT")
	router.Run(port)
}
