package frontend

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	ai "luillyfe.com/ai/semanticSearch"
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
	Query string `form:"query" binding:"required"`
}

// It handles the search query and returns a semantic similarity search result.
func postQueries(ctx *gin.Context) {
	// Will Hold the request data
	var request Request

	// Decode query text
	if ctx.ShouldBind(&request) == nil {
		log.Printf("Query received: %s", request.Query)
		processQuery(request.Query)
	}

	// Send back a success message
	ctx.JSON(200, gin.H{"message": "success!"})
}

func processQuery(query string) {
	// Get embedding for the current query
	featureVector := ai.GetEmbedding(query)

	// Query the index on vector search
	closestNeighbors := ai.Search(featureVector)

	fmt.Println(closestNeighbors)
	// Get results from query
	// return ai.GetResults(closestNeighbors)
}
