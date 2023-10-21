package ai

import (
	"context"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type semanticSearch struct {
	client *aiplatform.MatchClient
}

func NewSemanticSearch(ctx context.Context) *semanticSearch {
	client, err := aiplatform.NewMatchClient(ctx, option.WithEndpoint("102531040.us-central1-145252452137.vdb.vertexai.goog"))
	if err != nil {
		panic(err)
	}

	return &semanticSearch{client: client}
}

// Given an feature vector, it finds the most nearest neighbor in the previously set Dataset
func (s *semanticSearch) Query(ctx context.Context, embedding []float32) []*aiplatformpb.FindNeighborsResponse_NearestNeighbors {
	indexEndpoint := os.Getenv("GCP_INDEX_ENDPOINT")
	deployedIndexId := os.Getenv("GCP_INDEX_ID")

	queries := []*aiplatformpb.FindNeighborsRequest_Query{
		{Datapoint: &aiplatformpb.IndexDatapoint{
			DatapointId:   uuid.NewString(),
			FeatureVector: embedding,
		}},
	}

	request := &aiplatformpb.FindNeighborsRequest{
		IndexEndpoint:   indexEndpoint,
		DeployedIndexId: deployedIndexId,
		Queries:         queries,
	}
	response, _ := s.client.FindNeighbors(ctx, request)

	return response.GetNearestNeighbors()
}

func ToFloat32(s []interface{}) []float32 {
	result := make([]float32, len(s))
	for i, v := range s {
		elem := float32(v.(float64))
		result[i] = elem
	}

	return result
}
