package suggester

import (
	"github.com/jesusfar/meli.price.suggester/meli"
	"github.com/jesusfar/meli.price.suggester/mock"
	"github.com/stretchr/testify/assert"
	"log"
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

func TestSuggester_fetchItemsByCategory(t *testing.T) {
	categoryId := "MLA1050"
	mockServer := httptest.NewServer(http.HandlerFunc(mock.SearchItemsMock))
	defer mockServer.Close()

	os.Setenv("MELI_ENDPOINT", mockServer.URL)

	suggester := NewSuggester()

	suggester.fetchItemsByCategory(meli.SITE_MLA, categoryId)

	assert.Equal(t, true, directoryExists(DATA_SET_PATH+"/"+categoryId))
}

func TestSuggester_Train(t *testing.T) {
	suggester := NewSuggester()
	suggester.Train()
}

func TestSuggester_LoadDataTrained(t *testing.T) {
	suggester := NewSuggester()
	suggester.LoadDataTrained()
}

func TestSuggester_Suggest(t *testing.T) {
	suggester := NewSuggester()
	price, err := suggester.Suggest("MLA1050")
	log.Println(price)
	log.Println(err)
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
