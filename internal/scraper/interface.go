package scraper

import "github.com/gocolly/colly"


// ScrapedData can be any type of scraped data
type ScrapedData interface{}

// Scraper defines what every scraper must implement
type Scraper interface {
	GetName() string
	GetAllowedDomains() []string
	ConfigureCollector(c *colly.Collector) error
	Scrape(url string) ([]ScrapedData, error)
} 