package suggester

import (
	"github.com/jesusfar/meli.price.suggestor/meli"
	"log"
	"io/ioutil"
	"os"
	"fmt"
	"encoding/json"
	"sync"
	"errors"
)

const (
	FETCH_DATA_SET string = "fetchDataSet"
	TRAIN_MODEL string = "train"
	PREDICT string = "predict"
	DATA_SET_PATH = "./dataset/"
	DATA_TRAINED_PATH = "./datatrained/"
	DATA_TRAINED_FILE_PATH = DATA_TRAINED_PATH + "datatrained.json"
)

type DataTrained struct {
	sync.RWMutex
	data map[string]TrainedCategory
}

type TrainedCategory struct {
	Max       float64
	Suggested float64
	Min       float64
	Sum       float64
	Total     float64
}

type Suggester struct {
	meliClient meli.MeliClient
}

func NewSuggester() *Suggester {
	meliClient := meli.NewMeliHttpClient()
	//TODO remove this
	meliClient.SetEndpoint("http://localhost:3000")

	suggester := Suggester{
		meliClient: meliClient,
	}
	return &suggester
}

func (s *Suggester) FetchDataSet(site string)  {

	// Create folder if not exists
	createFolder(DATA_SET_PATH)

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

func (s *Suggester) Predict(categoryId string) (float64, error) {
	var dataTrained map[string]TrainedCategory
	var suggested float64

	dataTrainedFile, err := ioutil.ReadFile(DATA_TRAINED_FILE_PATH)

	if err != nil {
		return suggested, err
	}

	err = json.Unmarshal(dataTrainedFile, &dataTrained)

	if err != nil {
		return suggested, err
	}

	result, ok := dataTrained[categoryId]

	if ok {
		suggested = result.Suggested
		return suggested, nil
	} else {
		err = errors.New("CategoryId not found")
		return suggested, err
	}
}

func (s *Suggester) Train()  {

	wgItemProducer := &sync.WaitGroup{}
	wgItemConsumer := &sync.WaitGroup{}

	dataTrained := &DataTrained{data: make(map[string]TrainedCategory)}

	outPutItemChannel := make(chan *meli.SearchItem, 20)

	// Read dataSet path
	dataSetFolder, err :=ioutil.ReadDir(DATA_SET_PATH)

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


	dataTrainedForSave, err := json.Marshal(dataTrained.data)

	createFolder(DATA_TRAINED_PATH)

	ioutil.WriteFile(DATA_TRAINED_FILE_PATH, dataTrainedForSave, 0777)

	log.Println("[Train] Train finished")
}

func (s *Suggester) fetchItemsByCategory(site string, categoryId string)  {
	query := "category="+categoryId
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

func saveDataSet(searchItems []meli.SearchItem, categoryId string, index int)  {

	itemJson, _ := json.Marshal(searchItems)

	fileDest := fmt.Sprintf("%s/%s/%s-%d.json", DATA_SET_PATH, categoryId, categoryId, index)
	err := ioutil.WriteFile(fileDest, itemJson, 0777)
	if err != nil {
		log.Println(err)
	}
}

func createFolder(path string)  {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}


func trainModel(dataTrained *DataTrained, outPutItemChannel <-chan *meli.SearchItem, wg *sync.WaitGroup) {

	var dataTrain TrainedCategory

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
			suggested := sum/total

			dataTrain = TrainedCategory{
				Max: max,
				Min: min,
				Sum: sum,
				Total: total,
				Suggested: suggested,
			}

		} else {
			dataTrain = TrainedCategory{
				Max: itemInfo.Price,
				Min: itemInfo.Price,
				Sum: itemInfo.Price,
				Total: 1,
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

	for index, item := range items  {
		log.Printf("[readItemFile] Sending index: %d  item: %s", index, item.Id)
		itemToSend := item
		outPutItemChannel <- &itemToSend
	}
}