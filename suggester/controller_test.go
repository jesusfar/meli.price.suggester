package suggester

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const CategoryIdTest string = "MLA1051"

var dataTrainedTest map[string]CategoryPriceTrained

func init() {
	// Prepare dataTrained
	dataTrainedTest = make(map[string]CategoryPriceTrained)
	dataTrainedTest[CategoryIdTest] = CategoryPriceTrained{
		Max:       100.0,
		Suggested: 90.0,
		Min:       60,
	}
}

func TestNewSuggesterCtrl(t *testing.T) {
	suggesterCtrl := NewSuggesterCtrl()
	assert.NotNil(t, suggesterCtrl)
}

func TestSuggesterCtrl_SuggestPriceByCategory(t *testing.T) {

	expectedResult := `{"max":100,"suggested":90,"min":60}`

	t.Log("Given a categoryId: ", CategoryIdTest, " /categories/{categoryId}/prices returns a Suggest prices. ")

	{
		// Set dataTrained in Suggester
		s := NewSuggester()
		s.SetInMemoryDataTrained(dataTrainedTest)

		// Make a new Suggester Controller or handlers
		ctrl := SuggesterCtrl{Suggester: s}

		// Setup Gin test mode
		gin.SetMode(gin.TestMode)

		router := gin.New()

		router.GET("/categories/:categoryId/prices", ctrl.SuggestPriceByCategory)

		url := fmt.Sprintf("/categories/%s/prices", CategoryIdTest)

		// Send request
		req, _ := http.NewRequest("GET", url, nil)

		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// Assert expected result
		assert.Equal(t, resp.Body.String(), expectedResult)
	}
}

func BenchmarkSuggesterCtrl_SuggestPriceByCategory(b *testing.B) {

	b.ResetTimer()

	// Set dataTrained in Suggester
	s := NewSuggester()
	s.SetInMemoryDataTrained(dataTrainedTest)

	// Make a new Suggester Controller or handlers
	ctrl := SuggesterCtrl{Suggester: s}

	// Setup Gin test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.GET("/categories/:categoryId/prices", ctrl.SuggestPriceByCategory)

	url := fmt.Sprintf("/categories/%s/prices", CategoryIdTest)

	for i := 0; i < b.N; i++ {
		// Send request
		req, _ := http.NewRequest("GET", url, nil)

		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
	}
}
