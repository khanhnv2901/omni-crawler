package ecommerce

import (
	"github.com/gocolly/colly"
	"github.com/khanhnv2901/omni-crawler/internal/scraper"
	"github.com/khanhnv2901/omni-crawler/internal/types"
)

// Product represents a scraped product
type Product struct {
	Url   string `json:"url"`
	Image string `json:"image"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

// EcommerceScraper scrapers ecommerce sites
type EcommerceScraper struct {
	*scraper.BaseScraper           // Embbed struct (inheritance-like)
	products             []Product // Slice to store products
}

// NewScraper creates a new ecommerce scraper
func NewScraper(config *types.ScraperConfig) scraper.Scraper {
	return &EcommerceScraper{
		BaseScraper: scraper.NewBaseScraper(config),
		products:    make([]Product, 0), // empty slice
	}
}

// ConfigureCollector sets up scraping rules
func (e *EcommerceScraper) ConfigureCollector(c *colly.Collector) error {
	// Define what to do when we find a product
	c.OnHTML(e.Config.Selectors["product_item"], func(el *colly.HTMLElement) {
		product := Product{
			Url: el.ChildAttr(e.Config.Selectors["url"], "href"),
			Image: el.ChildAttr(e.Config.Selectors["image"], "src"),
			Name: el.ChildText(e.Config.Selectors["name"]),
			Price: el.ChildText(e.Config.Selectors["price"]),
		}
		e.products = append(e.products, product)
	})
	return nil
}

// Scrape performs the actual scraping
func (e *EcommerceScraper) Scrape(url string) ([]scraper.ScrapedData, error) {
	// Reset products
	e.products = make([]Product, 0)

	// Create and configure collector
	c := e.CreateCollector()
	if err := e.ConfigureCollector(c); err != nil {
		return nil, err
	}

	// Start scraping
	if err := c.Visit(url); err != nil {
		return nil, err
	}

	// Convert to interface slice
	data := make([]scraper.ScrapedData, len(e.products))
	for i, product := range e.products {
		data[i] = product
	}

	return data, nil
}
