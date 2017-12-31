package main

import (
	"os"
	"log"
	"github.com/jesusfar/meli.price.suggestor/suggestor"
	"github.com/jesusfar/meli.price.suggestor/meli"
	"fmt"
)

func main() {

	args := os.Args[1:]

	s := suggestor.NewSuggestor()

	switch args[0] {
	case suggestor.FETCH_DATA_SET:
		log.Println("[Suggestor] Fetch dataset")
		s.FetchDataSet(meli.SITE_MLA)
	default:
		log.Println("[Suggestor] Error action not defined.")
	}

	var input string
	fmt.Scanln(&input)
}
