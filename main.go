package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/peterzernia/lets-fork/restaurant"
	"github.com/peterzernia/lets-fork/utils"
	"github.com/peterzernia/lets-fork/websocket"
)

func main() {
	rdb, err := utils.InitRDB()
	if err != nil {
		log.Println(err)
	}

	pong, err := rdb.Ping().Result()
	log.Println(pong, err)

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

	router.Use(static.Serve("/", static.LocalFile("./public", true)))
	router.GET("", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	port := ":" + os.Getenv("PORT")
	router.Run(port)
}
