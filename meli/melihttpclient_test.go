package meli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMeliHttpClient(t *testing.T) {
	client := NewMeliHttpClient()

	assert.NotNil(t, client)
}

func TestMeliHttpClient_GetCategories_Returns_Error_Site_Validation(t *testing.T) {
	client := NewMeliHttpClient()
	_, err := client.GetCategories("")
	assert.EqualError(t, err, "[MeliHttpClientErr] Error description: Site param mustn't be empty.")
}
