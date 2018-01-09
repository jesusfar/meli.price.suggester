package meli

import (
	"encoding/json"
	"fmt"
	"github.com/mercadolibre/golang-restclient/rest"
	"log"
	"net/http"
	"os"
)

const MELI_API_ENDPOINT = "https://api.mercadolibre.com"

type MeliHttpClient struct {
	endpoint string
}

func NewMeliHttpClient() *MeliHttpClient {

	endpoint := os.Getenv("MELI_ENDPOINT")

	if endpoint == "" {
		endpoint = MELI_API_ENDPOINT
	}

	log.Println("[MeliHttpClient] Set endpoint: " + endpoint)

	client := MeliHttpClient{endpoint: endpoint}

	return &client
}

func (m *MeliHttpClient) SetEndpoint(endpoint string) {
	m.endpoint = endpoint
}

func (m *MeliHttpClient) GetCategories(site string) ([]Category, error) {

	var categories []Category

	if site == "" {
		err := MeliClientErr{Message: "Site param mustn't be empty."}
		return nil, err
	}

	url := fmt.Sprintf("%s/sites/%s/categories", m.endpoint, site)

	res := rest.Get(url)

	if isSuccess(res) {
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

	if isSuccess(res) {
		err := json.Unmarshal(res.Bytes(), &searchItems)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		return &searchItems, nil
	} else {
		err := MeliClientErr{Message: "Error searching items"}
		return nil, err
	}

}

func isSuccess(res *rest.Response) bool {
	if res.Response != nil && res.StatusCode == http.StatusOK {
		return true
	}
	log.Println("[MeliHttpClient] Response is nil or status code is not success.")
	log.Println(res)
	return false
}

// MeliClientErr
type MeliClientErr struct {
	Message string
	Code    string
}

func (e MeliClientErr) Error() string {
	return fmt.Sprintf("[MeliHttpClientErr] Error description: %s", e.Message)
}
