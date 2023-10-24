package ai

import (
	"context"
	"fmt"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

type predictionClient struct {
	endpoint string
	client   *aiplatform.PredictionClient
}

// TODO: Better handling errors
func NewPredictionClient(ctx context.Context, vertexAIEndpoint string) *predictionClient {
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

	return &predictionClient{endpoint: endpoint, client: predictionServiceClient}
}

// TODO: There is no need for a channel. whether make it concurrent or remove channel.
func (s *predictionClient) Predict(ctx context.Context, predictions chan []*structpb.Value, instances []*structpb.Value) {
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
		Instances:  instances,
		Parameters: parameters,
	})
	if err != nil {
		panic(err)
	}

	predictions <- response.Predictions
}

func (*predictionClient) BuildInstance(content string) *structpb.Value {
	// Adding text embeddings parameters
	// https://cloud.google.com/vertex-ai/docs/generative-ai/model-reference/text-embeddings#request_body
	instances, err := structpb.NewValue(map[string]interface{}{"content": content})
	if err != nil {
		panic(err)
	}

	return instances
}

// https: //cloud.google.com/vertex-ai/docs/generative-ai/model-reference/text-embeddings#response_body
func GetVectors(predictions []*structpb.Value) []InputData {
	vectorEmbeddingsChan := make(chan InputData)
	vectorEmbeddings := make([]InputData, 0)

	for _, p := range predictions {
		go func(p *structpb.Value) {
			for _, v := range p.GetStructValue().Fields {
				embedding := v.GetStructValue().Fields["values"].GetListValue().AsSlice()
				vectorEmbeddingsChan <- InputData{Id: uuid.New(), Embedding: embedding}
			}
		}(p)
	}

	for i := 0; i < len(predictions); i++ {
		vector := <-vectorEmbeddingsChan
		vectorEmbeddings = append(vectorEmbeddings, vector)
	}

	return vectorEmbeddings
}
