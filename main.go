package main

import (
	"log"
	"os"

	"gopkg-analyzer/src/processor"
)

func main() {
	f, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = processor.AnalyzeCode(f)
	if err != nil {
		log.Fatal(err)
	}
}
