package meli

const(
	SITE_MLA string = "MLA"
)

// MeliClient defines base interface operation
type MeliClient interface {
	GetCategories(site string) ([]Category, error)
	SearchItems(query string)
}


