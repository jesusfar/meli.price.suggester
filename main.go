package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jesusfar/meli.price.suggester/meli"
	"github.com/jesusfar/meli.price.suggester/suggester"
	"log"
	"os"
	"github.com/jesusfar/meli.price.suggester/api"
)

func printHelp() {
	fmt.Println(`

MeliPriceSugesster is a tool for predict a price by category

Usage: meliPriceSugesster <command>

Commands:

  fetch            Fetch data set by categories
  train	           Train the data set
  suggest          Suggest a price given a category
  serve            Serve a http service
  help             Help Meli Price Suggester

Examples:
  meliPriceSugesster fetchDataSet
  meliPriceSugesster train
  meliPriceSugesster serve 3000
  meliPriceSugesster predict MLA70400

	`)
}

func serve() {

	s := api.NewSuggesterCtrl()

	r := gin.Default()

	r.GET("/categories/:categoryId/prices", s.SuggestPriceByCategory)

	r.Run(":8080")
}

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
	case suggester.SUGGEST:
		if len(args) >= 2 {
			categoryId := args[1]
			priceSuggested, err := s.Suggest(categoryId)
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

	case suggester.SERVE:
		serve()
	default:
		printHelp()
	}

	var input string
	fmt.Scanln(&input)
}
