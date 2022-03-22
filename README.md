# `Prodact`

`Prodact` is a small Go command line program that fetches Oda's product catalogue with a user-defined rate limit and simple filtering options.

`Prodact` was built with Go 1.17.



## First: discarded alternative solutions

### 1. Use product categories

- Oda's products are divided into categories, where each category might have multiple child caterogries.
- Oda's top-level categories are exposed as JSON under this endpoint: `https://oda.com/api/v1/app-components/products/`.

**JSON snippet:**

```json
{
  "title": "Varer",
  "blocks": [
    {
      "id": "link-list-product-listings",
      "component": "link-list",
      "links": [
        {
          "title": "Dine varer",
          "target": {
            "method": "push",
            "title": "Dine varer",
            "uri": "https://oda.com/no/products/yours/"
          },
          "image": {
            "uri": "https://oda.com/static/img/app/icon-star-32@3x.35782c960fe7.png"
          }
        }
      ]
    },
    {
      "id": "link-list-product-categories",
      "component": "link-list",
      "links": [
        {
          "title": "Frukt og grønt",
          "target": {
            "method": "push",
            "title": "Frukt og grønt",
            "uri": "https://oda.com/no/categories/20-frukt-og-gront/"
          },
          "image": {
            "uri": "https://oda.com/static/product_categorization/img/svg/pear.1a8a8b0a8c2e.svg"
          }
        }
      ],
      "title": "Kategorier"
    }
  ]
}
```

Each category has its own URL (`https://oda.com/api/v1/app-components/products/<ID>`) with child categories for each category living under the `links` key. If a category's `links` array is empty, it means it has no child categories, and its human URL can be scraped for products.

One way of building a product catalogue would consist of two parts:
* first, traverse the JSON output from `app-components` to create a category tree:
  - process url strings to get category IDs, then get JSON definition of the category.
* then, process HTMS that lives under each bottom-category's link: find product name, price, currency, etc., and present that data to the user.

This solution would easily be extended to allow users to limit their search to one or more categories. It was discarded due to time constrains, as it required two very different ways of gathering data: one unmarshalling JSON, the other processing HTML.



### 2. Exploit regex-based routes

* [Oda has an API](https://github.com/kolonialno/api-docs) that is not generally available (but as we saw above, some endpoints are public). Under that API, one can get a product's information by visiting: `https://oda.com/no/api/v1/products/<ID>/`.
* On the human-facing website though, the URLs for each product contain not only the ID, but also the name of the product, like so: `https://oda.com/no/products/26541-pink-lady-epler-italia/`.
* Visiting `https://oda.com/no/products/26541/` returns a 404 ☹️
* But visiting `https://oda.com/no/products/26541-foobar/`... redirects to `https://oda.com/no/products/26541-pink-lady-epler-italia/`

Thus, another way of building a catalogue would be to:
* range over a set number of such URLs,
* then, process the HTML to find product name, price, currency, etc. live, and present that data to the user.

This solution was discarded because product IDs are unique, and there is no reasonable programmatic way to figue out how high the number should be (*Example:* Pink Lady apples have an ID of 26541, and belong to the 20 -> 21 -> 513 category).


### 3. Build a traditional website crawler

A traditional website crawler processes a website's main page and stores every url it can find in a queue. It then visits each url in that queue and does the whole process over again, until all links are in the inventory.

This solution was discarded as potentially expensive given the original purpose of the tool (a product catalogue).


## Chosen solution

`Prodact` leverages Oda's search API that's live under `https://oda.com/api/v1/search/`. 

The search API takes two (known) parameters: `q` for query and `page` for paginating the results.

## Usage

You can use this program without building the Go executable by running the following in the command line:

`$ go run cmd/prodact/main.go`

Or you can use `make build` to build the executable. You can then run it like so:

`$ ./cmd/prodact/prodact`

Options:

```
  -query string
    	put your search query here; leave empty for full inventory (default "")
  -rate duration
    	defines the pause between calls to the online store in seconds (default 1)
  -store string
    	which grocery store do you want to parse (default "oda")
  -unavailable
    	use this flag if you want to see unavailable products (default false)
```

## Data format

Currently, the program returns the products in a simple, `;`-separated list (not `,`, because some items contain `,` in their name):

```
product ID; product full name; product price and currency; product availability
35;Møllerens Hvetemel Siktet;9.90NOK;false
32487;Natural Basics Refill Hand Wash White Tea & Verbena 500ml;32.00NOK;false
1855;Nestlé Cheerios Multi Frokostblanding;25.40NOK;true
26431;Telys Maxi 10 Timer;51.90NOK;true
16984;Semper Minibaguetter Glutenfri;50.50NOK;true
23349;Hipp Combiotik Pulver 2 Fra 6 - 12 mnd;107.00NOK;true
28891;Dr. Oetker Rustica Pizza Double Pepperoni;69.90NOK;false
24889;Domstein Ørretfilet med Skinn Oppdrett, Fersk;54.24NOK;true
27540;Sørlandskjøtt Servelat;24.90NOK;true
```

## Testing

Writing test functions was skipped due to time constraints.
