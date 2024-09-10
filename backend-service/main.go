package main

import (
	"backend-service/routes"
	utils "backend-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Connect Redis cache here
	utils.ConnectToRedis()

	// Test route to check functionality
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Health is good!",
		})
	})

	// Backend routes registration
	routes.GameRoutes(r)
	//routes.AuthRoutes(r)
	//...

	// Start the server
	r.Run()
}
