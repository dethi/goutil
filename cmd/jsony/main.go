package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	var payload interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&payload); err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(b)
	os.Stdout.WriteString("\n")
}
