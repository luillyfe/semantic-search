package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"google.golang.org/protobuf/types/known/structpb"
	ai "luillyfe.com/ai/semanticSearch"
)

// Given a JSON files it reads from
func ReadJSON(fileName string, Decoder interface{}, textEmdeddings chan interface{}) {
	// Open the JSON file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode the JSON file
	err = json.NewDecoder(file).Decode(&Decoder)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Pass the text to the channel
	textEmdeddings <- Decoder
}

// Given a JSON Lines files it reads from
func ReadJSONL(fileName string, Decoder ai.AIDataset, linesChan chan []interface{}) {
	// Open the JSON Lines file
	lines, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// Decode the JSON Lines file
	data := make([]interface{}, 0)
	for _, line := range strings.Split(string(lines), "\n") {
		if err := json.Unmarshal([]byte(line), &Decoder); err != nil {
			panic(err)
		}
		data = append(data, Decoder)
	}

	// Pass the text to the channel
	linesChan <- data
}

func WriteJSONL(name string, vectors []ai.InputData) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	for _, v := range vectors {
		line, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString(string(line) + "\n")
		if err != nil {
			panic(err)
		}
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

// Given a Chan that holds JSONL lines it will write them to a JSONL file when they are ready
func WriteJSONLInBatches(
	name string,
	predictionsChan chan []*structpb.Value,
	buildVectors func(predictions []*structpb.Value) []ai.InputData,
	numReqs int) {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	// Create the file
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	// Loop over the channel and for each prediction Response
	// write a batch to the JSON file (JSONL format)
	for i := 0; i < numReqs; i++ {
		predictions := <-predictionsChan
		wg.Add(1)
		go func(vectors []ai.InputData) {
			defer wg.Done()
			for _, v := range vectors {
				line, err := json.Marshal(v)
				if err != nil {
					panic(err)
				}

				mutex.Lock()
				_, err = f.WriteString(string(line) + "\n")
				mutex.Unlock()
				if err != nil {
					panic(err)
				}
			}
		}(buildVectors(predictions))
	}

	wg.Wait()

	// Close the file
	err = f.Close()
	if err != nil {
		panic(err)
	}

}
