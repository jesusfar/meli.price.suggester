package meli

import (
	"encoding/json"
	"errors"
	"github.com/jesusfar/meli.price.suggester/mock"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const checkMark = "\u2713"

func init() {
	log.Println("Init test")
}

func TestNewMeliHttpClient(t *testing.T) {
	client := NewMeliHttpClient()
	assert.NotNil(t, client)
}

func TestMeliHttpClient_GetCategories(t *testing.T) {

	var categories []Category
	var categoriesResult []Category

	// Read categories
	file, _ := mock.ReadFileOfCategories()
	json.Unmarshal(file, &categoriesResult)

	// Set testCases to test
	var testCases = []struct {
		messageTest string
		inputSite   string
		result      []Category
		err         error
	}{
		{
			messageTest: "Given a empty site GetCategories returns error.",
			inputSite:   "",
			result:      categories,
			err:         errors.New("[MeliHttpClientErr] Error description: Site param mustn't be empty."),
		},
		{
			messageTest: "Given a site MLA GetCategories returns a collection of categories.",
			inputSite:   "MLA",
			result:      categoriesResult,
		},
	}

	// Run Mock server
	server := httptest.NewServer(http.HandlerFunc(mock.GetCategoriesMock))
	defer server.Close()

	client := NewMeliHttpClient()
	client.SetEndpoint(server.URL)

	t.Log("TestCase GetCategories")
	{
		for _, testCase := range testCases {
			t.Log(testCase.messageTest, checkMark)

			result, err := client.GetCategories(testCase.inputSite)

			if err != nil {
				assert.Equal(t, err.Error(), testCase.err.Error())
			} else {
				assert.Equal(t, testCase.result, result)
			}
		}
	}
}

func TestMeliHttpClient_SearchItems(t *testing.T) {

	// Run Mock server
	server := httptest.NewServer(http.HandlerFunc(mock.SearchItemsMock))
	defer server.Close()

	client := NewMeliHttpClient()
	client.SetEndpoint(server.URL)

	t.Log("TestCase SearchItems")
	{
		result, _ := client.SearchItems(SITE_MLA, "category=MLA1051", 0, 50)

		t.Log("SearchItems return SearchItemResult", checkMark)
		assert.NotNil(t, result)
		assert.IsType(t, &SearchItemsResult{}, result)
	}

}

func TestMeliHttpClient_SetEndpoint(t *testing.T) {
	endpoint := "http://localhost:3000"

	client := NewMeliHttpClient()
	client.SetEndpoint(endpoint)

	t.Log("SetEnpoint ", endpoint, checkMark)
	assert.Equal(t, endpoint, client.GetEndpoint())
}
