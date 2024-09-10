package routes

import (
	controllers "backend-service/controllers/game"

	"github.com/gin-gonic/gin"
)

// GameRoutes defines routes related to game operations
func GameRoutes(r *gin.Engine) {
	gameGroup := r.Group("/games")
	{
		gameGroup.POST("/live", controllers.CreateLiveGameRoom)
		gameGroup.POST("/local", controllers.CreateLocalGameRoom)
	}
}
