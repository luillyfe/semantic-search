package ai

import (
	"context"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func GetEmbedding(text string) []InputData {
	ctx := context.Background()

	// AI platform regional endpoint
	endpoint := "us-central1-aiplatform.googleapis.com:443"
	// Get prediction client
	predictionClient := NewPredictionClient(ctx, endpoint)

	// Format query for request to the prediction API
	instances := predictionClient.BuildInstance(text)

	instanceChan := make(chan []*structpb.Value)
	predictionClient.Predict(ctx, instanceChan, []*structpb.Value{instances})

	predictResponse := <-instanceChan
	//
	return GetVectors(predictResponse)
}

func Search(embedding []InputData) []*aiplatformpb.FindNeighborsResponse_NearestNeighbors {
	ctx := context.Background()
	vectorSearchClient := NewSemanticSearch(ctx)

	return vectorSearchClient.Query(ctx, ToFloat32(embedding[0].Embedding))
}
