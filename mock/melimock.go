package mock

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func SearchItemsMock(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	file, _ := ReadFileSearchItems()

	categories := string(file[:])

	fmt.Fprintln(w, categories)
}

func GetCategoriesMock(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	file, _ := ReadFileOfCategories()

	categories := string(file[:])

	fmt.Fprintln(w, categories)
}

func ReadFileOfCategories() ([]byte, error) {
	file, err := ioutil.ReadFile("./../mock/Get-Categories-MLA.json")

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return file, nil
}

func ReadFileSearchItems() ([]byte, error) {
	file, err := ioutil.ReadFile("./../mock/Search-By-Caterogry-MLA1051.json")

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return file, nil
}
