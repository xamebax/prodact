package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// StoreURL defines which online store we want to scrape
const StoreURL = "https://oda.com"

// RateLimit defines the pause between calls to the online store
const RateLimit = 1 * time.Second

var log = logrus.New()

func main() {
	log.Infof("🧡 starting prODAct with a rate limit of %s...", RateLimit)

	BuildProductCatalogue()
	log.Info("✨ my job is done, exiting cleanly...")
}

// BuildProductCatalogue builds the online store's product catalogue.
func BuildProductCatalogue() {
	log.Info("starting to build product catalogue...")
	filename := fmt.Sprintf("fixtures/products-%s.txt", time.Now().Round(time.Minute))
	file, _ := os.Create(filename)
	defer file.Close()
	productCatalogue := Catalogue{}
	// We now from manual inspection that the search API returns 40 results per page.
	// We'll subtract this number from the totalHits that we got the first time
	// to know when to stop paging, as the number of products (totalHits) is not a constant.
	// resultsPerPage := 40
	// at the moment of writing, an empty search gives the value of 7377;
	// so we can expect the results to be around 7377/40 ~= 185 pages
	// totalHits := 0

	for page := 1; page < 10; page++ {
		searchResults := SearchResults{}
		url := fmt.Sprintf("%s/api/v1/search/?page=%d&q=", StoreURL, page)
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("error while getting url: %s", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("cannot read response budy: %s", err)
		}
		err = json.Unmarshal(body, &searchResults)
		if err != nil {
			log.Errorf("cannot unmarshal search results: %s", err)
		}
		productCatalogue.Products = append(productCatalogue.Products, searchResults.Products...)
		log.Infof("found %d products on page #%d...", len(searchResults.Products), page)
		time.Sleep(RateLimit)
	}
	log.Infof("finished compiling product catalogue, writing it to file %s", file.Name())
	for i, product := range productCatalogue.Products {
		_, err := file.WriteString(fmt.Sprintf("%d. %d, %s, %s%s, %t\n", i, product.ID, product.Name, product.GrossPrice, product.Currency, product.Availability.IsAvailable))
		if err != nil {
			log.Errorf("cannot write results to file: %s", err)
		}
	}
}

// Catalogue is where all of our data is collected. A Catalogue contains two
// fields, one for storing products, and one for storing metadata for
// categorizing the products
type Catalogue struct {
	Categories []Category
	Products   []Product
}

// To save time typing, structs below were first created using https://mholt.github.io/json-to-go/

// Product contains metadata of a single product item.
type Product struct {
	ID           int          `json:"id"`
	FullName     string       `json:"full_name"`
	Brand        string       `json:"brand"`
	BrandID      int          `json:"brand_id"`
	Name         string       `json:"name"`
	NameExtra    string       `json:"name_extra"`
	FrontURL     string       `json:"front_url"`
	AbsoluteURL  string       `json:"absolute_url"`
	GrossPrice   string       `json:"gross_price"`
	Currency     string       `json:"currency"`
	Discount     interface{}  `json:"discount"`
	Promotion    interface{}  `json:"promotion"`
	Availability Availability `json:"availability"`
}

// Availability answers the question of whether an item is available to buy.
type Availability struct {
	IsAvailable bool `json:"is_available"`
}

// SearchResults is a top-level struct into which we unmarshal a page of search
// results at a time.
type SearchResults struct {
	Attributes Attributes    `json:"attributes"`
	Products   []Product     `json:"products"`
	Categories []interface{} `json:"categories"`
}

// Attributes contains information of the number of total hits
// a search result returns.
type Attributes struct {
	TotalHits int `json:"total_hits"`
}

// Category contains data on a single product category and its children, if there are any.
/*
An example JSON response with a single category:
{
  "title": "Bakeri og brød",
  "target": {
    "method": "push",
      "title": "Bakeri og brød",
      "uri": "https://oda.com/no/categories/1135-bakeri-og-brod/"
  },
  "image": {
    "uri": "https://oda.com/static/product_categorization/img/svg/bread.f1b088ffddf0.svg"
  }
},
*/
type Category struct {
	Title    string `json:"title,omitempty"`
	ID       string
	Children []Child // We assume that Oda's links array is a list of child categories

}

// Child would contain a single child category. If its links are empty,
// we can assume it is a bottom-level category and is safe to scrape for products.
type Child struct {
	// TODO
}

// BuildCategories visits Oda's API products endpoint:
// https://oda.com/api/v1/app-components/products/ and first parses the
// top-level categories JSON. It then traverses through the links in each
// category and builds a master struct
func BuildCategories() {
	log.Info("starting to build category catalogue...")
	// store.BuildCatalogue()

}
