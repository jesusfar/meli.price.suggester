package meli

const (
	SITE_MLA string = "MLA"
)

// MeliClient defines base interface operation
type MeliClient interface {
	GetCategories(site string) ([]Category, error)
	SearchItems(site string, query string, offset int, limit int) (*SearchItemsResult, error)
}

type Category struct {
	Id   string
	Name string
}

type SearchItemsResult struct {
	SiteId  string       `json:"site_id"`
	Paging  PageInfo     `json:"paging"`
	Results []SearchItem `json:"results"`
}

type PageInfo struct {
	Total          int `json:"total"`
	Offset         int `json:"offset"`
	Limit          int `json:"limit"`
	PrimaryResults int `json:"primary_results"`
}

type SearchItem struct {
	Id         string  `json:"id"`
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	Currency   string  `json:"currency_id"`
	CategoryId string  `json:"category_id"`
}
