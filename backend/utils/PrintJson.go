package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	"log"
)

func LogJson[T any](input T) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", jsonData)
}
