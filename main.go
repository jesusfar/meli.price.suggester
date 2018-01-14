package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jesusfar/meli.price.suggester/meli"
	"github.com/jesusfar/meli.price.suggester/suggester"
	"os"
)

func printHelp() {
	fmt.Println(`

MeliPriceSugesster is a tool for predict a price by category

Usage: meliPriceSugesster <command>

Commands:

  fetch            Fetch data set by categories
  train	           Train the data set
  suggest          Suggest a price given a category
  clean            Clean data set and data trained folders
  serve            Serve a http service
  help             Help Meli Price Suggester

Examples:
  meliPriceSugesster fetch
  meliPriceSugesster train
  meliPriceSugesster serve 3000
  meliPriceSugesster predict MLA70400

	`)
}

func serve() {

	s := suggester.NewSuggesterCtrl()

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
		s.FetchDataSet(meli.SITE_MLA)
	case suggester.TRAIN_MODEL:
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
	case suggester.CLEAN:
		s.Clean()
	default:
		printHelp()
	}

	var input string
	fmt.Scanln(&input)
}
