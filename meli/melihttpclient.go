package meli

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/mercadolibre/golang-restclient/rest"
	"net/http"
)

const MELI_API_ENDPOINT  = "https://api.mercadolibre.com/"

type MeliHttpClient struct {

}

func NewMeliHttpClient() *MeliHttpClient {

	client := MeliHttpClient{}

	return &client
}

func (m *MeliHttpClient) GetCategories(site string) ([]Category, error) {

	var categories []Category

	if site == "" {
		err := MeliClientErr{Message: "Site param mustn't be empty."}
		return nil, err
	}

	url := fmt.Sprintf("%s/sites/%s/categories", MELI_API_ENDPOINT, site)


	res := rest.Get(url)

	if res.Response != nil && res.StatusCode == http.StatusOK {
		err := json.Unmarshal(res.Bytes(), &categories)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		return categories, nil
	} else {
		err := MeliClientErr{Message: "Error fetching for categories"}
		return nil, err
	}
}

func (m * MeliHttpClient) SearchItems(query string)  {

}

type MeliClientErr struct {
	Message string
}

func (e MeliClientErr) Error() string {
	return fmt.Sprintf("[MeliHttpClientErr] Error description: %s", e.Message)
}
