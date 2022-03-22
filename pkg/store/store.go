package store

import (
	"fmt"
	"time"
)

// BuildProductCatalogue is the package's public API function.
// It calls store-specific code depending on which store we want to scrape.
func BuildProductCatalogue(products chan<- Product, errors chan<- error, storeName string, rateLimit time.Duration, searchQuery string) {
	switch storeName {
	case "Oda", "oda", "ODA":
		buildODAProductCatalogue(products, errors, rateLimit, searchQuery)
	default:
		panic(fmt.Sprintf("prODAct doesn't support scraping %s yet", storeName))
	}
}
