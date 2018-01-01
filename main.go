package main

import (
	"os"
	"log"
	"github.com/jesusfar/meli.price.suggestor/meli"
	"fmt"
	"github.com/jesusfar/meli.price.suggestor/suggester"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	s := suggester.NewSuggester()

	switch args[0] {
	case suggester.FETCH_DATA_SET:
		log.Println("[Suggestor] Fetch dataset")
		s.FetchDataSet(meli.SITE_MLA)
	case suggester.TRAIN_MODEL:
		log.Println("[Suggestor] Train model")
		s.Train()
	case suggester.PREDICT:
		if len(args) >= 2 {
			categoryId := args[1]
			result, err := s.Predict(categoryId)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Printf("Price suggested: %f", result)
		} else {
			printHelp()
		}

	default:
		printHelp()
	}

	var input string
	fmt.Scanln(&input)
}

func printHelp() {
	fmt.Println("Meli Price Suggestor help.")
}
