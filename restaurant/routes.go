package restaurant

import (
	"github.com/gin-gonic/gin"
)

// InitializeRoutes initializes routes for the App
func InitializeRoutes(r *gin.RouterGroup) {
	restaurant := r.Group("/restaurants")

	restaurant.GET("/:id", handleGet)
}
