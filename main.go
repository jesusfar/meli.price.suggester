package main

import (
	"fmt"
	"github.com/jesusfar/meli.price.suggester/meli"
	"github.com/jesusfar/meli.price.suggester/suggester"
	"log"
	"os"
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
		log.Println("[suggester] Fetch dataset")
		s.FetchDataSet(meli.SITE_MLA)
	case suggester.TRAIN_MODEL:
		log.Println("[suggester] Train model")
		s.Train()
	case suggester.PREDICT:
		if len(args) >= 2 {
			categoryId := args[1]
			priceSuggested, err := s.Predict(categoryId)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Printf("For category: %s  Price suggested: %f , Min: %f, Max: %f",
				categoryId,
				priceSuggested.Suggested,
				priceSuggested.Min,
				priceSuggested.Max)
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
	fmt.Println(`

MeliPriceSugesster is a tool for predict a price by category

Usage: meliPriceSugesster <command>

Commands:

  fetchDataSet     Fetch data set by categories
  train	           Train the data set
  predict          Predict a price given a category
  serve            Serve a http service
  help             Help Meli Price Suggester

Examples:
  meliPriceSugesster fetchDataSet
  meliPriceSugesster train
  meliPriceSugesster serve 3000
  meliPriceSugesster predict MLA70400

	`)
}
