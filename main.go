package main

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	ai "luillyfe.com/ai/semanticSearch"
	"luillyfe.com/ai/utils"
)

func main() {
	ctx := context.Background()
	// AI platform regional endpoint
	endpoint := "us-central1-aiplatform.googleapis.com:443"
	// Get a prediction client
	search := ai.NeWSemanticSearch(ctx, endpoint)

	// The dataset demonstrates the use of the Text Embedding API with a vector database.
	// gs://cloud-samples-data/vertex-ai/dataset-management/datasets/bert_finetuning/wide_and_deep_trainer_container_tests_input.jsonl
	fileName := "wide_and_deep_trainer_container_tests_input.jsonl"
	linesChan := make(chan []interface{})
	go utils.ReadJSONL(fileName, ai.AIDataset{}, linesChan)

	// Build the text to embed #limit to one line to ease Results interpretation
	lines := <-linesChan
	dataFrame := ai.NewDataFrame(search.BuildInstance, lines)

	// Get Prediction Response
	predictionsChan := make(chan []*structpb.Value)

	// The prediction API for AutoML models has a restriction of 5 Instances per request.
	reqInstancesLimit := 5
	numOfRequests := 0
	for i := 0; i <= len(dataFrame); i += reqInstancesLimit {
		numOfRequests++

		end := i + reqInstancesLimit
		if end > len(dataFrame) {
			end = len(dataFrame)
		}

		go func(dataInBatch []*structpb.Value) {
			search.Predict(ctx, predictionsChan, dataInBatch)
		}(dataFrame[i:end])
		// TODO: Come up with a better rate limiting algorithm
		time.Sleep(12 * time.Second)
	}

	// Writing vector embeddings in batches to JSONL file
	utils.WriteJSONLInBatches(
		"vectorEmbeddings.json",
		predictionsChan,
		ai.GetVectors,
		// Number of issued requests to the Prediction API
		numOfRequests,
	)

	// NewJobServiceClient
}
