package ai

import "github.com/google/uuid"

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

type InputData struct {
	Id        uuid.UUID     `json:"id"`
	Embedding []interface{} `json:"embedding"`
}
