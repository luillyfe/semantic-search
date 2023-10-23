package frontend

import "github.com/gin-gonic/gin"

func Run() {
	// Creating a router
	router := gin.Default()

	// API Definition
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "success!"})
	})

	// Run
	router.Run("localhost:8080")
}
