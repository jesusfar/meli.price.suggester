package suggestor

import (
	"github.com/jesusfar/meli.price.suggestor/meli"
	"log"
	"github.com/gin-gonic/gin/json"
	"io/ioutil"
	"os"
	"fmt"
)

const (
	FETCH_DATA_SET string = "fetchDataSet"
	TRAIN_MODEL string = "train"
	PREDICT string = "predict"
)

type Suggestor struct {
	meliClient meli.MeliClient
}

func NewSuggestor() *Suggestor  {
	meliClient := meli.NewMeliHttpClient()
	//TODO remove this
	meliClient.SetEndpoint("http://localhost:3000")

	suggestor := Suggestor{
		meliClient: meliClient,
	}
	return &suggestor
}

func (s *Suggestor) FetchDataSet(site string)  {


	// Fetch categories for site MLA
	categories, err := s.meliClient.GetCategories(site)

	if err != nil {
		log.Println("[FetchDataSet] Error fetching categories")
		return
	}

	// Foreach category we need to search items related
	for _, category := range categories {
		log.Println("[FetchDataSet] Fetching items for category: " + category.Id)

		go s.fetchItemsByCategory(site, category.Id)
	}
}

func (s *Suggestor) Train()  {
	
}

func (s *Suggestor) fetchItemsByCategory(site string, categoryId string)  {
	query := "category="+categoryId
	offset := 0
	limit := 50

	createFolder("./dataset/"+categoryId)

	searchResult, err := s.meliClient.SearchItems(site, query, offset, limit)

	if err != nil {
		log.Println("[searchItemsByCategory] Error searching items.")
		return
	}

	// Save first DataSet
	saveDataSet(searchResult.Results, categoryId, 0)

	sampleSize := int(CalcSampleSize(searchResult.Paging.Total))

	log.Printf("[fetchItemsByCategory] Sample Size: %d", sampleSize)

	offset = limit

	for offset < sampleSize {
		// Search Item offset
		searchResult, err := s.meliClient.SearchItems(site, query, offset, limit)

		if err != nil {
			log.Println("[searchItemsByCategory] Error searching items.")
			return
		}

		saveDataSet(searchResult.Results, categoryId, offset)

		offset += limit
	}
}

func saveDataSet(searchItems []meli.SearchItem, categoryId string, index int)  {

	itemJson, _ := json.Marshal(searchItems)

	fileDest := fmt.Sprintf("dataset/%s/%s-%d.json", categoryId, categoryId, index)
	err := ioutil.WriteFile(fileDest, itemJson, 0777)
	log.Println(err)
}

func createFolder(path string)  {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}