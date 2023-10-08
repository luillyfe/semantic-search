package main

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
	ai "luillyfe.com/ai/semanticSearch"
)

func main() {
	ctx := context.Background()

	// AI platform regional endpoint
	endpoint := "us-central1-aiplatform.googleapis.com:443"
	search := ai.NeWSemanticSearch(ctx, endpoint)

	content := "What is life?"
	predictions := make(chan []*structpb.Value)
	go search.Predict(ctx, predictions, content)

	fmt.Println(<-predictions)
}
