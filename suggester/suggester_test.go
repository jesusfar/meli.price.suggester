package suggester

import (
	"github.com/jesusfar/meli.price.suggester/meli"
	"github.com/jesusfar/meli.price.suggester/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const checkMark = "\u2713"

func TestNewSuggester(t *testing.T) {

	suggester := NewSuggester()
	t.Log("NewSuggester returns a Suggester pointer.", checkMark)
	assert.NotNil(t, suggester)
	assert.IsType(t, &Suggester{}, suggester)

}

func TestSuggester_FetchDataSet(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(mock.GetCategoriesMock))
	defer mockServer.Close()

	os.Setenv("MELI_ENDPOINT", mockServer.URL)

	suggester := NewSuggester()

	suggester.FetchDataSet(meli.SITE_MLA)
}

func TestSuggester_SetInMemoryDataTrained(t *testing.T) {

	// Prepare data trained for test
	dataTrainedTest = make(map[string]CategoryPriceTrained)
	dataTrainedTest[CategoryIdTest] = CategoryPriceTrained{
		Max:       100.0,
		Suggested: 90.0,
		Min:       60,
	}

	s := NewSuggester()

	s.SetInMemoryDataTrained(dataTrainedTest)

	if assert.Equal(t, dataTrainedTest, s.GetInMemoryDataTrained().data) {
		t.Log("Given a data trained, set suggester with in memory data trained.", checkMark)
	}

}

func TestNewSuggester_fetchItemsBySystematicRandomSampling(t *testing.T) {
	categoryId := "MLA1050"
	mockServer := httptest.NewServer(http.HandlerFunc(mock.SearchItemsMock))
	defer mockServer.Close()

	os.Setenv("MELI_ENDPOINT", mockServer.URL)

	suggester := NewSuggester()

	suggester.FetchItemsBySystematicRandomSampling(meli.SITE_MLA, categoryId)

	assert.Equal(t, true, directoryExists(DATA_SET_PATH+"/"+categoryId))
}

func TestSuggester_Train(t *testing.T) {
	suggester := NewSuggester()
	suggester.Train()
	assert.Equal(t, true, directoryExists(DATA_TRAINED_PATH))
}

func TestSuggester_LoadDataTrained(t *testing.T) {
	suggester := NewSuggester()
	suggester.LoadDataTrained()
}

func TestSuggester_Suggest(t *testing.T) {
	// Prepare data trained for test
	dataTrainedTest = make(map[string]CategoryPriceTrained)
	dataTrainedTest[CategoryIdTest] = CategoryPriceTrained{
		Max:       100.0,
		Suggested: 90.0,
		Min:       60,
	}

	expectedResult := CategoryPriceSuggested{Max: 100, Suggested: 90, Min: 60}

	s := NewSuggester()

	s.SetInMemoryDataTrained(dataTrainedTest)

	price, err := s.Suggest(CategoryIdTest)

	assert.IsType(t, CategoryPriceSuggested{}, price)
	assert.Nil(t, err)
	if assert.Equal(t, expectedResult, price) {
		t.Log("Given a CategoryId: ", CategoryIdTest, " suggester returns: ", expectedResult, checkMark)
	}
}

func TestSuggester_Clean(t *testing.T) {
	suggester := NewSuggester()
	suggester.Clean()

	assert.Equal(t, false, directoryExists(DATA_SET_PATH))
	assert.Equal(t, false, directoryExists(DATA_TRAINED_PATH))
}

func directoryExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}
	return true
}
