# Semantic Search with Vertex AI

This project provides a step-by-step guide on how to build a semantic search engine using the Vertex AI APIs. Semantic search is a more sophisticated mechanism to find relevant content to a search rather than traditionally keyword-based. It takes relevant information that it may be not be present in the query. It does so by understanding the context of the input text that the user types on the search box like: user query history, location of user input, among others.

This project uses the following Vertex AI APIs:

**Text Embeddings API**: To generate text embeddings for the search index and the query text.

**Vector Search API**: To perform semantic search on the text embeddings.

# Getting Started

To get started, you will need to create a Google Cloud project. Once you have created a project, you will need to enable the Vertex AI API.

Building the Search Index
To build the search index, you will need to first generate text embeddings for the documents that you want to include in the index. You can use the Text Embeddings API to generate text embeddings.

Once you have generated text embeddings for the documents, you will need to upload the embeddings to Vector Search. You can use the Vector Search API to upload the embeddings.

Performing Semantic Search
To perform semantic search, you will need to send a query text to the Vector Search API. The Vector Search API will return a list of documents that are semantically similar to the query text.

Example
The following code shows how to perform semantic search using the Vertex AI APIs:

Go

```go
import (
    aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
)
```

# Get the Text Embeddings API client.

```go
predictionServiceClient, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(vertexAIEndpoint))
```

# Generate text embeddings for a given text.

```go
response, err := s.client.Predict(ctx, &aiplatformpb.PredictRequest{
		Endpoint:   s.endpoint,
		Instances:  []*structpb.Value{instances},
		Parameters: parameters,
	})
```

# Train a model using AutoML

```go
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
```

# Get the Vector Search API client.

```go
indexEndpoint := os.Getenv("GCP_INDEX_ENDPOINT")
deployedIndexId := os.Getenv("GCP_INDEX_ID")

client, err := aiplatform.NewMatchClient(ctx, option.WithEndpoint("102531040.us-central1-145252452137.vdb.vertexai.goog"))
if err != nil {
	panic(err)
}
```

...

# Perform semantic search on the text embeddings.

...

# Querying the model

```go
// Query the model
queries := []*aiplatformpb.FindNeighborsRequest_Query{
	{Datapoint: &aiplatformpb.IndexDatapoint{
		DatapointId:   uuid.New(),
		FeatureVector: embedding,
	}},
}

request := &aiplatformpb.FindNeighborsRequest{
	IndexEndpoint:   indexEndpoint,
	DeployedIndexId: deployedIndexId,
	Queries:         queries,
}
// Find all your neighbors
response, _ := client.FindNeighbors(ctx, request)

// Get the closest neighbor to your feature vector (Your query)
response.GetNearestNeighbors()
```

...

# Conclusion

This project provides a step-by-step guide on how to build a Semantic Search engine using the Vertex AI APIs. Semantic search is a powerful tool that can be used to improve the search experience for your users.
