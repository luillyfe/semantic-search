package ai

import (
	"context"
	"fmt"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
)

// Given an feature vector, it finds the most nearest neighbor in the previously set Dataset
func Query(ctx context.Context, embedding []float32) {
	indexEndpoint := os.Getenv("GCP_INDEX_ENDPOINT")
	deployedIndexId := os.Getenv("GCP_INDEX_ID")

	client, err := aiplatform.NewMatchClient(ctx, option.WithEndpoint("102531040.us-central1-145252452137.vdb.vertexai.goog"))
	if err != nil {
		panic(err)
	}

	queries := []*aiplatformpb.FindNeighborsRequest_Query{
		{Datapoint: &aiplatformpb.IndexDatapoint{
			DatapointId:   "0",
			FeatureVector: embedding,
		}},
	}

	request := &aiplatformpb.FindNeighborsRequest{
		IndexEndpoint:   indexEndpoint,
		DeployedIndexId: deployedIndexId,
		Queries:         queries,
	}
	response, _ := client.FindNeighbors(ctx, request)

	fmt.Println(response.GetNearestNeighbors())
}

func ToFloat32(s []interface{}) []float32 {
	result := make([]float32, len(s))
	for i, v := range s {
		elem := float32(v.(float64))
		result[i] = elem
	}

	return result
}
