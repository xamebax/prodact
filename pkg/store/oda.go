package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// SearchResults is a top-level struct into which we unmarshal a page of search
// results at a time.
type SearchResults struct {
	Attributes Attributes    `json:"attributes"`
	Products   []Product     `json:"products"`
	Categories []interface{} `json:"categories"`
}

// Product contains metadata of a single product item.
type Product struct {
	ID           int          `json:"id"`
	FullName     string       `json:"full_name"`
	Name         string       `json:"name"`
	FrontURL     string       `json:"front_url"`
	GrossPrice   string       `json:"gross_price"`
	Currency     string       `json:"currency"`
	Availability Availability `json:"availability"`
}

// Availability answers the question of whether an item is available to buy.
type Availability struct {
	IsAvailable bool `json:"is_available"`
}

// Attributes contains information of the number of total hits
// a search result returns.
type Attributes struct {
	TotalHits int `json:"total_hits"`
}

func buildODAProductCatalogue(products chan<- Product, errors chan<- error, rateLimit time.Duration, searchQuery string) error {
	defer close(products)

	for page := 1; ; page++ {
		searchResults := SearchResults{}
		url := fmt.Sprintf("https://oda.com/api/v1/search/?page=%d&q=%s", page, searchQuery)
		resp, err := http.Get(url)
		if err != nil {
			errors <- fmt.Errorf("error while getting url: %s", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errors <- fmt.Errorf("cannot read response body: %s", err)
		}

		err = json.Unmarshal(body, &searchResults)
		if err != nil {
			errors <- fmt.Errorf("cannot unmarshal search results: %s", err)
		}
		if len(searchResults.Products) == 0 {
			break
		}
		for _, product := range searchResults.Products {
			products <- product
		}
		time.Sleep(rateLimit)
	}

	return nil
}
