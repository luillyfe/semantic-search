Semantic Search with Vertex AI

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

# Generate text embeddings for the query text.

```go
response, err := s.client.Predict(ctx, &aiplatformpb.PredictRequest{
		Endpoint:   s.endpoint,
		Instances:  []*structpb.Value{instances},
		Parameters: parameters,
	})
```

# Get the Vector Search API client.

...

# Perform semantic search on the text embeddings.

...

# Print the search results.

...

# Conclusion

This project provides a step-by-step guide on how to build a Semantic Search engine using the Vertex AI APIs. Semantic search is a powerful tool that can be used to improve the search experience for your users.
