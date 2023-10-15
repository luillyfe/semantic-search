package ai

import (
	"google.golang.org/protobuf/types/known/structpb"
)

// DataFrame like structure
func NewDataFrame(buildInstance func(line string) *structpb.Value, lines []interface{}) []*structpb.Value {
	dataFrameChan := make(chan *structpb.Value, len(lines))
	dataFrame := make([]*structpb.Value, 0)

	for _, line := range lines {
		go func(line interface{}) {
			dataLine := line.(AIDataset)
			dataFrameChan <- buildInstance(dataLine.TextContent)

		}(line)
	}

	for {
		select {
		case dataLine := <-dataFrameChan:
			dataFrame = append(dataFrame, dataLine)
		default:
			if len(dataFrame) == len(lines) {
				return dataFrame
			}
		}
	}
}

func Filter(slice []*structpb.Value, filter func(v *structpb.Value) bool) []*structpb.Value {
	filterIn := make([]*structpb.Value, 0)
	for _, value := range slice {
		if filter(value) {
			filterIn = append(filterIn, value)
		}
	}

	return filterIn
}
