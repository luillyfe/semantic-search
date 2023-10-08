package ai

import (
	"context"
	"fmt"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

type semanticSearch struct {
	endpoint string
	client   *aiplatform.PredictionClient
}

// TODO: Better handling errors
func NeWSemanticSearch(ctx context.Context, vertexAIEndpoint string) *semanticSearch {
	// Instantiates a client
	predictionServiceClient, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(vertexAIEndpoint))
	if err != nil {
		panic(err)
	}

	projectId := os.Getenv("GCP_PROJECT_ID")
	region := os.Getenv("GCP_REGION")
	publisher := "google"
	model := "textembedding-gecko@001"

	// Configuring the parent resource
	endpoint := fmt.Sprintf(
		"projects/%s/locations/%s/publishers/%s/models/%s",
		projectId,
		region,
		publisher,
		model,
	)

	return &semanticSearch{endpoint: endpoint, client: predictionServiceClient}
}

func (s *semanticSearch) Predict(ctx context.Context, predictions chan []*structpb.Value, content string) {
	instances, err := structpb.NewValue(map[string]interface{}{"content": content})
	if err != nil {
		panic(err)
	}

	parameters, err := structpb.NewValue(map[string]interface{}{
		"temperature":     0,
		"maxOutputTokens": 256,
		"topP":            0,
		"topK":            1,
	})
	if err != nil {
		panic(err)
	}

	response, err := s.client.Predict(ctx, &aiplatformpb.PredictRequest{
		Endpoint:   s.endpoint,
		Instances:  []*structpb.Value{instances},
		Parameters: parameters,
	})
	if err != nil {
		panic(err)
	}

	predictions <- response.Predictions
}
