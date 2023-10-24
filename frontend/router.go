package frontend

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Run() {
	// Creating a router
	router := gin.Default()

	// Queries API
	router.POST("/", postQueries)

	// Run
	router.Run("localhost:8080")
}

type Request struct {
	Query string `form:"query"`
}

// It handles the search query and returns a semantic similarity search result.
func postQueries(ctx *gin.Context) {
	// Will Hold the request data
	var request Request

	// Decode query text
	if ctx.ShouldBind(&request) == nil {
		log.Println(request.Query)
	}

	// Send back a success message
	ctx.JSON(200, gin.H{"message": "success!"})
}
