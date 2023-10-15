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

func (s *semanticSearch) Predict(ctx context.Context, predictions chan []*structpb.Value, instances []*structpb.Value) {
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

func (*semanticSearch) BuildInstance(content string) *structpb.Value {
	// Adding text embeddings parameters
	// https://cloud.google.com/vertex-ai/docs/generative-ai/model-reference/text-embeddings#request_body
	instances, err := structpb.NewValue(map[string]interface{}{"content": content})
	if err != nil {
		panic(err)
	}

	return instances
}

type InputData struct {
	Id        int           `json:"id"`
	Embedding []interface{} `json:"embedding"`
}

// https: //cloud.google.com/vertex-ai/docs/generative-ai/model-reference/text-embeddings#response_body
func GetVectors(predictions []*structpb.Value) []InputData {
	vectorEmbeddingsChan := make(chan InputData)
	vectorEmbeddings := make([]InputData, 0)

	// TODO: Generate id
	for id, p := range predictions {
		go func(id int, p *structpb.Value) {
			for _, v := range p.GetStructValue().Fields {
				embedding := v.GetStructValue().Fields["values"].GetListValue().AsSlice()
				vectorEmbeddingsChan <- InputData{Id: id, Embedding: embedding}
			}
		}(id, p)
	}

	for {
		select {
		case vector := <-vectorEmbeddingsChan:
			vectorEmbeddings = append(vectorEmbeddings, vector)
		default:
			if len(vectorEmbeddings) == len(predictions) {
				return vectorEmbeddings
			}
		}
	}
}

type AIDataset struct {
	TextContent              string                   `json:"textContent"`
	ClassificationAnnotation ClassificationAnnotation `json:"classificationAnnotation"`
	DataItemResourceLabels   interface{}              `json:"dataItemResourceLabels"`
	// Embedding                interface{}              `json:"embedding"`
}

type ClassificationAnnotation struct {
	DisplayName              string      `json:"displayName"`
	AnnotationResourceLabels interface{} `json:"annotationResourceLabels"`
}
