package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func OpenJSONFile(fileName string, Decoder interface{}, textEmdeddings chan interface{}) {
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
