package websocket

import (
	"github.com/gin-gonic/gin"
)

// InitializeRoutes initializes routes for the App
func InitializeRoutes(r *gin.RouterGroup, h *Hub) {
	websocket := r.Group("/")

	websocket.GET("ws", func(c *gin.Context) {
		serve(h, c.Writer, c.Request)
	})
}
