package suggester

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuggesterCtrl struct {
	Suggester *Suggester
}

func NewSuggesterCtrl() *SuggesterCtrl {
	s := &SuggesterCtrl{
		Suggester: NewSuggester(),
	}

	return s
}

func (s *SuggesterCtrl) SuggestPriceByCategory(c *gin.Context) {
	categoryId := c.Param("categoryId")

	// Validate param
	if len(categoryId) == 0 {
		c.JSON(http.StatusBadRequest, ApiErr{Message: "CategoryId param is empty."})
		return
	}

	// Suggest prices for category
	result, err := s.Suggester.Suggest(categoryId)

	if err != nil {
		c.JSON(http.StatusNotFound, ApiErr{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

type ApiErr struct {
	Message string
}

func (e ApiErr) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}
