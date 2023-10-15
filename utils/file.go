package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
