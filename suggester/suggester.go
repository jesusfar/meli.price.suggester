package suggester

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jesusfar/meli.price.suggester/meli"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	FETCH_DATA_SET         string = "fetch"
	TRAIN_MODEL            string = "train"
	SUGGEST                string = "suggest"
	SERVE                  string = "serve"
	DATA_SET_PATH                 = "./dataset/"
	DATA_TRAINED_PATH             = "./datatrained/"
	DATA_TRAINED_FILE_PATH        = DATA_TRAINED_PATH + "datatrained.json"
)

type DataTrained struct {
	sync.RWMutex
	data map[string]CategoryPriceTrained
}

type CategoryPriceTrained struct {
	Max       float64
	Suggested float64
	Min       float64
	Sum       float64
	Total     float64
}

type CategoryPriceSuggested struct {
	Max       float64 `json:"max"`
	Suggested float64 `json:"suggested"`
	Min       float64 `json:"min"`
}

type Suggester struct {
	meliClient          meli.MeliClient
	inMemoryDataTrained *DataTrained
}

// NewSuggester returns a suggester for category price.
func NewSuggester() *Suggester {
	meliClient := meli.NewMeliHttpClient()

	suggester := &Suggester{
		meliClient: meliClient,
	}

	// LoadDataTrained if exists data trained file
	suggester.LoadDataTrained()

	return suggester
}

// FetchDataSet fetches items from Meli and save data in dataset folder
func (s *Suggester) FetchDataSet(site string) {

	// Create folder if not exists
	createFolder(DATA_SET_PATH)

	// Fetch categories for site
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

// Suggest a price for categoryId
func (s *Suggester) Suggest(categoryId string) (CategoryPriceSuggested, error) {
	var suggested CategoryPriceSuggested

	// Try to load data trained.
	if s.inMemoryDataTrained == nil {
		err := s.LoadDataTrained()
		if err != nil {
			log.Println("[Predict] Error loading data trained.")
			return suggested, err
		}
	}

	result, ok := s.inMemoryDataTrained.data[categoryId]

	if ok {
		suggested.Max = result.Max
		suggested.Suggested = result.Suggested
		suggested.Min = result.Min
		return suggested, nil
	} else {
		err := errors.New(fmt.Sprintf("Category: %s not found.", categoryId))
		return suggested, err
	}
}

// Train reads the dataSet and prepare the model to predict the price by categoryID
func (s *Suggester) Train() {

	wgItemProducer := &sync.WaitGroup{}
	wgItemConsumer := &sync.WaitGroup{}

	dataTrained := &DataTrained{data: make(map[string]CategoryPriceTrained)}

	outPutItemChannel := make(chan *meli.SearchItem, 20)

	// Read dataSet path
	dataSetFolder, err := ioutil.ReadDir(DATA_SET_PATH)

	if err != nil {
		log.Println("[Train] Error reading ./dataset/ folder")
	}

	for _, file := range dataSetFolder {
		if file.IsDir() {
			categoryId := file.Name()
			log.Println("[Train] Starting train dataset for category: " + categoryId)

			wgItemProducer.Add(1)
			go readItemFilesForCategory(categoryId, outPutItemChannel, wgItemProducer)

			wgItemConsumer.Add(1)
			go trainModel(dataTrained, outPutItemChannel, wgItemConsumer)
		}
	}

	log.Println("[Train] Waiting to finish")

	wgItemProducer.Wait()
	close(outPutItemChannel)

	wgItemConsumer.Wait()

	dataTrainedForSave, _ := json.Marshal(dataTrained.data)

	createFolder(DATA_TRAINED_PATH)

	err = ioutil.WriteFile(DATA_TRAINED_FILE_PATH, dataTrainedForSave, 0777)

	if err != nil {
		log.Println("[Train] Error writing data trained.")
		log.Println(err)
	}

	// Reset dataTrained in Suggester
	s.inMemoryDataTrained = nil

	log.Println("[Train] Train finished")
}

// LoadDataTrained loads data trained from file if exist and keep in memory.
func (s *Suggester) LoadDataTrained() error {
	var dataTrained map[string]CategoryPriceTrained

	dataTrainedFile, err := ioutil.ReadFile(DATA_TRAINED_FILE_PATH)

	if err != nil {
		log.Printf("[LoadDataTrained][Notice] Data trained file: %s does not exist.", DATA_TRAINED_FILE_PATH)
		return err
	}

	err = json.Unmarshal(dataTrainedFile, &dataTrained)

	if err != nil {
		log.Printf("[LoadDataTrained][Notice] Error Unmarshal file: %d ", DATA_TRAINED_FILE_PATH)
		log.Println(err)
		return err
	}

	inMemoryDataTrained := &DataTrained{data: dataTrained}

	// Load dataTrained
	s.inMemoryDataTrained = inMemoryDataTrained

	log.Print("[LoadDataTrained][Notice]  Data trained load [OK]")

	return nil
}

func (s *Suggester) fetchItemsByCategory(site string, categoryId string) {
	query := "category=" + categoryId
	offset := 0
	limit := 50

	createFolder(DATA_SET_PATH + categoryId)

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

func saveDataSet(searchItems []meli.SearchItem, categoryId string, index int) {

	itemJson, _ := json.Marshal(searchItems)

	fileDest := fmt.Sprintf("%s/%s/%s-%d.json", DATA_SET_PATH, categoryId, categoryId, index)
	err := ioutil.WriteFile(fileDest, itemJson, 0777)
	if err != nil {
		log.Println(err)
	}
}

func createFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}

func trainModel(dataTrained *DataTrained, outPutItemChannel <-chan *meli.SearchItem, wg *sync.WaitGroup) {

	var dataTrain CategoryPriceTrained

	// Iterate while outPutItemChannel is open
	for itemInfo := range outPutItemChannel {
		log.Printf("[trainModel] Item: %s", itemInfo.Id)

		categoryId := itemInfo.CategoryId

		dataTrained.RLock()
		value, exists := dataTrained.data[categoryId]
		dataTrained.RUnlock()

		if exists {
			max := value.Max
			min := value.Min

			if itemInfo.Price > value.Max {
				max = itemInfo.Price
			} else if itemInfo.Price < value.Min {
				min = itemInfo.Price
			}

			sum := value.Sum + itemInfo.Price
			total := value.Total + 1
			suggested := sum / total

			dataTrain = CategoryPriceTrained{
				Max:       max,
				Min:       min,
				Sum:       sum,
				Total:     total,
				Suggested: suggested,
			}

		} else {
			dataTrain = CategoryPriceTrained{
				Max:       itemInfo.Price,
				Min:       itemInfo.Price,
				Sum:       itemInfo.Price,
				Total:     1,
				Suggested: itemInfo.Price,
			}
		}

		dataTrained.Lock()
		dataTrained.data[categoryId] = dataTrain
		dataTrained.Unlock()
	}

	wg.Done()
	log.Println("[trainModel] Done.")
}

func readItemFilesForCategory(categoryId string, outPutItemChannel chan<- *meli.SearchItem, wg *sync.WaitGroup) {

	categoryDataSetPath := DATA_SET_PATH + categoryId

	log.Printf("[readCategory:%s] Reading dataset from: %s", categoryId, categoryDataSetPath)

	dataSetFiles, err := ioutil.ReadDir(categoryDataSetPath)

	if err != nil {
		log.Printf("[readCategory:%s] Error reading dataset: %s", categoryId, categoryDataSetPath)
		return
	}

	for _, file := range dataSetFiles {
		if !file.IsDir() {
			filePath := categoryDataSetPath + "/" + file.Name()

			readItemFileForCategory(categoryId, filePath, outPutItemChannel)
		}
	}

	wg.Done()
}

func readItemFileForCategory(categoryId string, filePath string, outPutItemChannel chan<- *meli.SearchItem) {

	var items []meli.SearchItem

	log.Printf("[readItemCategory:%s] Reading file: %s", categoryId, filePath)

	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Printf("[readItemCategory:%s] Error reading file: %s", categoryId, filePath)
	}

	err = json.Unmarshal(file, &items)

	if err != nil {
		log.Printf("[readItemCategory:%s] Error Unmarshal file: %s", categoryId, filePath)
	}

	for index, item := range items {
		log.Printf("[readItemFile] Sending index: %d  item: %s", index, item.Id)
		itemToSend := item
		outPutItemChannel <- &itemToSend
	}
}
