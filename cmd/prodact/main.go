package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xamebax/prodact/pkg/store"
)

var (
	searchQuery     = flag.String("query", "", "put your search query here; leave empty for full inventory")
	rateLimit       = flag.Duration("rate", 1*time.Second, "defines the pause between calls to the online store in seconds; pass 0 to disable rate limiting")
	storeName       = flag.String("store", "oda", "which grocery store do you want to parse; leave empty for Oda")
	showUnavailable = flag.Bool("unavailable", false, "use this flag if you want to include unavailable products")
	log             = logrus.New()
)

func main() {
	flag.Parse()

	log.Info("🧡 starting prODAct...")
	log.Infof("scraping %s; rate limit of 1 req/%s; showing unavailable products: %t; query: %s",
		*storeName, *rateLimit, *showUnavailable, *searchQuery)

	products := make(chan store.Product, 40)
	errors := make(chan error, 40)

	go store.BuildProductCatalogue(products, errors, *storeName, *rateLimit, *searchQuery)

	for {
		select {
		case product, ok := <-products:
			if !ok {
				log.Info("✨ my job is done, exiting cleanly...")
				return
			}
			if *showUnavailable || product.Availability.IsAvailable {
				// We could use a CSV-exporting library here, too.
				_, err := fmt.Fprintf(os.Stdout, "%d;%s;%s%s;%t\n",
					product.ID, product.FullName,
					product.GrossPrice, product.Currency,
					product.Availability.IsAvailable,
				)
				if err != nil {
					log.Panicf("can't show products because of an unrecoverrable error: %s", err)
				}
			}
		case err := <-errors:
			log.Error(err)
		}
	}
}
