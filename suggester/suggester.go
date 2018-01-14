package suggester

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jesusfar/meli.price.suggester/meli"
	"github.com/jesusfar/meli.price.suggester/util"
	"io/ioutil"
	"os"
	"sync"
)

const (
	FETCH_DATA_SET         string = "fetch"
	TRAIN_MODEL            string = "train"
	SUGGEST                string = "suggest"
	SERVE                  string = "serve"
	CLEAN                  string = "clean"
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
	logger              *util.Logger
}

// NewSuggester returns a suggester for category price.
func NewSuggester() *Suggester {
	meliClient := meli.NewMeliHttpClient()

	suggester := &Suggester{
		meliClient: meliClient,
		logger:     util.NewLogger(),
	}

	// LoadDataTrained if exists data trained file
	suggester.LoadDataTrained()

	return suggester
}

// FetchDataSet fetches items from Meli and save data in dataset folder
func (s *Suggester) FetchDataSet(site string) {

	s.logger.Info("[FetchDataSet] Fetching data set ...")

	// Create folder if not exists
	createFolder(DATA_SET_PATH)

	// Fetch categories for site
	categories, err := s.meliClient.GetCategories(site)

	if err != nil {
		s.logger.Info("[FetchDataSet] Error fetching categories. Please see in DEBUG mode")
		s.logger.Debug(err)
		return
	}

	// Foreach category we need to search items related
	for _, category := range categories {
		s.logger.Debug("[FetchDataSet] Fetching items for category: " + category.Id)
		//go s.fetchItemsByCategory(site, category.Id)
		go s.fetchItemsBySystematicRandomSampling(site, category.Id)
	}

	s.logger.Info("[FetchDataSet] Fetching done.")
}

// Suggest a price for categoryId
func (s *Suggester) Suggest(categoryId string) (CategoryPriceSuggested, error) {
	var suggested CategoryPriceSuggested

	// Try to load data trained.
	if s.inMemoryDataTrained == nil {
		err := s.LoadDataTrained()
		if err != nil {
			s.logger.Warning("[Predict] Error loading data trained.")
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
		s.logger.Warning("[Train] Error reading ./dataset/ folder")
	}

	for _, file := range dataSetFolder {
		if file.IsDir() {
			categoryId := file.Name()
			s.logger.Debug("[Train] Starting train dataset for category: " + categoryId)

			wgItemProducer.Add(1)
			go s.readItemFilesForCategory(categoryId, outPutItemChannel, wgItemProducer)

			wgItemConsumer.Add(1)
			go s.trainModel(dataTrained, outPutItemChannel, wgItemConsumer)
		}
	}

	s.logger.Info("[Train] Waiting to finish")

	wgItemProducer.Wait()
	close(outPutItemChannel)

	wgItemConsumer.Wait()

	dataTrainedForSave, _ := json.Marshal(dataTrained.data)

	createFolder(DATA_TRAINED_PATH)

	err = ioutil.WriteFile(DATA_TRAINED_FILE_PATH, dataTrainedForSave, 0777)

	if err != nil {
		s.logger.Warning("[Train] Error writing data trained.")
		s.logger.Debug(err)
	}

	// Reset dataTrained in Suggester
	s.inMemoryDataTrained = nil

	s.logger.Info("[Train] Train finished")
}

// LoadDataTrained loads data trained from file if exist and keep in memory.
func (s *Suggester) LoadDataTrained() error {
	var dataTrained map[string]CategoryPriceTrained

	dataTrainedFile, err := ioutil.ReadFile(DATA_TRAINED_FILE_PATH)

	if err != nil {
		s.logger.Warning("[LoadDataTrained][Notice] Data trained file: %s does not exist.", DATA_TRAINED_FILE_PATH)
		return err
	}

	err = json.Unmarshal(dataTrainedFile, &dataTrained)

	if err != nil {
		s.logger.Warning(fmt.Sprintf("[LoadDataTrained][Notice] Error Unmarshal file: %d ", DATA_TRAINED_FILE_PATH))
		s.logger.Debug(err)
		return err
	}

	s.SetInMemoryDataTrained(dataTrained)

	s.logger.Info("[LoadDataTrained][Notice]  Data trained load [OK]")

	return nil
}

func (s *Suggester) SetInMemoryDataTrained(data map[string]CategoryPriceTrained) {
	s.logger.Info("[SetInMemoryDataTrained] Set in memory data trained.")
	s.inMemoryDataTrained = &DataTrained{data: data}
}

func (s *Suggester) GetInMemoryDataTrained() *DataTrained {
	return s.inMemoryDataTrained
}

func (s *Suggester) fetchItemsBySystematicRandomSampling(site string, categoryId string) {

	query := "category=" + categoryId
	offset := 0
	limit := 50

	createFolder(DATA_SET_PATH + categoryId)

	searchResult, err := s.meliClient.SearchItems(site, query, offset, limit)

	if err != nil {
		s.logger.Warning("[fetchRandomItemsByCategory] Error searching items.")
		return
	}

	// Save first DataSet
	s.saveDataSet(searchResult.Results, categoryId, 0)

	// Fetch next items by Systematic Random Sampling

	// Get total sampling
	totalItems := searchResult.Paging.Total
	s.logger.Info(fmt.Sprintf("[fetchItemsByCategory][%s] Total Items: %d", categoryId, totalItems))

	// Get sample size
	sampleSize := util.CalcSampleSizeMethod2(totalItems)
	s.logger.Info(fmt.Sprintf("[fetchItemsByCategory][%s] Sample Size: %d", categoryId, sampleSize))

	// Calc P elements p = N / n where N is total items and n is sample size
	p := totalItems / sampleSize
	s.logger.Info(fmt.Sprintf("[fetchItemsByCategory][%s] Proportion of elements p: %d", categoryId, p))

	// Calc K, where offsetK is random offset to start.
	offsetK := util.GetRandomNumberFrom(p)
	s.logger.Info(fmt.Sprintf("[fetchItemsByCategory][%s] Initial offset: %d", categoryId, offsetK))

	i := 0
	nextOffsetK := 0
	for nextOffsetK < totalItems {
		i++
		nextOffsetK = offsetK + i*p
		s.logger.Debug(fmt.Sprintf("[fetchItemsByCategory][%s] Next offset: %d  index: %d", categoryId, nextOffsetK, i))

		searchResult, err := s.meliClient.SearchItems(site, query, nextOffsetK, limit)

		if err != nil {
			s.logger.Warning("[searchItemsByCategory] Error searching items.")
			s.logger.Debug(err)
			return
		}

		s.saveDataSet(searchResult.Results, categoryId, nextOffsetK)
	}
}

// Clean removes data set and data trained folders.
func (s *Suggester) Clean() {
	s.logger.Info("[Clean] Cleaning data..")
	s.cleanDirectory(DATA_SET_PATH)
	s.cleanDirectory(DATA_TRAINED_PATH)
	s.logger.Info("[Clean] Done.")
}

func (s *Suggester) cleanDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		s.logger.Warning(err)
	}
}

func (s *Suggester) saveDataSet(searchItems []meli.SearchItem, categoryId string, index int) {

	itemJson, _ := json.Marshal(searchItems)

	fileDest := fmt.Sprintf("%s/%s/%s-%d.json", DATA_SET_PATH, categoryId, categoryId, index)
	err := ioutil.WriteFile(fileDest, itemJson, 0777)
	if err != nil {
		s.logger.Warning("[saveDataSet] Error saving dataset.")
		s.logger.Debug(err)
	}
}

func createFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}

func (s *Suggester) trainModel(dataTrained *DataTrained, outPutItemChannel <-chan *meli.SearchItem, wg *sync.WaitGroup) {

	var dataTrain CategoryPriceTrained

	// Iterate while outPutItemChannel is open
	for itemInfo := range outPutItemChannel {
		s.logger.Debug(fmt.Sprintf("[trainModel] Item: %s", itemInfo.Id))

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
	s.logger.Debug("[trainModel] Done.")
}

func (s *Suggester) readItemFilesForCategory(categoryId string, outPutItemChannel chan<- *meli.SearchItem, wg *sync.WaitGroup) {

	categoryDataSetPath := DATA_SET_PATH + categoryId

	s.logger.Debug(fmt.Sprintf("[readCategory:%s] Reading dataset from: %s", categoryId, categoryDataSetPath))

	dataSetFiles, err := ioutil.ReadDir(categoryDataSetPath)

	if err != nil {
		s.logger.Warning("[readCategory:%s] Error reading dataset: %s", categoryId, categoryDataSetPath)
		s.logger.Debug(err)
		return
	}

	for _, file := range dataSetFiles {
		if !file.IsDir() {
			filePath := categoryDataSetPath + "/" + file.Name()

			s.readItemFileForCategory(categoryId, filePath, outPutItemChannel)
		}
	}

	wg.Done()
}

func (s *Suggester) readItemFileForCategory(categoryId string, filePath string, outPutItemChannel chan<- *meli.SearchItem) {

	var items []meli.SearchItem

	s.logger.Debug(fmt.Sprintf("[readItemCategory:%s] Reading file: %s", categoryId, filePath))

	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		s.logger.Warning(fmt.Sprintf("[readItemCategory:%s] Error reading file: %s", categoryId, filePath))
		s.logger.Debug(err)
	}

	err = json.Unmarshal(file, &items)

	if err != nil {
		s.logger.Warning(fmt.Sprintf("[readItemCategory:%s] Error Unmarshal file: %s", categoryId, filePath))
		s.logger.Debug(err)
	}

	for index, item := range items {
		s.logger.Debug(fmt.Sprintf("[readItemFile] Sending index: %d  item: %s", index, item.Id))
		itemToSend := item
		outPutItemChannel <- &itemToSend
	}
}
