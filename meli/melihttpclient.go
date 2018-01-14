package meli

import (
	"encoding/json"
	"fmt"
	"github.com/jesusfar/meli.price.suggester/util"
	"github.com/mercadolibre/golang-restclient/rest"
	"log"
	"net/http"
	"os"
)

const MELI_API_ENDPOINT = "https://api.mercadolibre.com"

type MeliHttpClient struct {
	endpoint string
	logger   *util.Logger
}

func NewMeliHttpClient() *MeliHttpClient {

	endpoint := os.Getenv("MELI_ENDPOINT")

	if endpoint == "" {
		endpoint = MELI_API_ENDPOINT
	}

	client := MeliHttpClient{
		endpoint: endpoint,
		logger:   util.NewLogger(),
	}

	return &client
}

func (m *MeliHttpClient) SetEndpoint(endpoint string) {
	m.endpoint = endpoint
}

func (m *MeliHttpClient) GetEndpoint() string {
	return m.endpoint
}

func (m *MeliHttpClient) GetCategories(site string) ([]Category, error) {

	var categories []Category

	if site == "" {
		err := MeliClientErr{Message: "Site param mustn't be empty."}
		return nil, err
	}

	url := fmt.Sprintf("%s/sites/%s/categories", m.endpoint, site)

	res := rest.Get(url)

	if m.isSuccess(res) {
		err := json.Unmarshal(res.Bytes(), &categories)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		return categories, nil
	} else {
		err := MeliClientErr{Message: "Error fetching for categories."}
		return nil, err
	}
}

func (m *MeliHttpClient) SearchItems(site string, query string, offset int, limit int) (*SearchItemsResult, error) {
	var searchItems SearchItemsResult

	url := fmt.Sprintf("%s/sites/%s/search?q=%s&offset=%v&limit=%v", m.endpoint, site, query, offset, limit)

	res := rest.Get(url)

	if m.isSuccess(res) {
		err := json.Unmarshal(res.Bytes(), &searchItems)

		if err != nil {
			m.logger.Debug(err)
			return nil, err
		}

		return &searchItems, nil
	} else {
		err := MeliClientErr{Message: "Error searching items"}
		return nil, err
	}

}

func (m *MeliHttpClient) isSuccess(res *rest.Response) bool {
	if res.Response != nil && res.StatusCode == http.StatusOK {
		return true
	}
	m.logger.Warning("[MeliHttpClient] Response is nil or status code is not success.")
	m.logger.Debug(res)
	return false
}

// MeliClientErr
type MeliClientErr struct {
	Message string
}

func (e MeliClientErr) Error() string {
	return fmt.Sprintf("[MeliHttpClientErr] Error description: %s", e.Message)
}
