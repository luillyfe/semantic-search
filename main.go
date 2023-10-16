package main

import (
	"context"

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
	dataFrame := ai.NewDataFrame(search.BuildInstance, lines[:1])

	// Get Prediction Response
	predictionsChan := make(chan []*structpb.Value)
	go search.Predict(ctx, predictionsChan, dataFrame)

	// Writing vector embeddings to file
	predictions := <-predictionsChan
	vectorEmbeddings := ai.GetVectors(predictions)
	// Vector search accepts a jsonl file but it does required to name it as .json
	utils.WriteJSONL("vectorEmbeddings.json", vectorEmbeddings)

	// NewJobServiceClient
}
