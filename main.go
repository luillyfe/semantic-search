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

	input := "What is Final Fantasy?"
	predictions := make(chan []*structpb.Value)
	go search.Predict(ctx, predictions, input)

	fmt.Println(<-predictions)
}
