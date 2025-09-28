package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/khanhnv2901/omni-crawler/internal/types"
)

// BaseScraper provides common functionality
type BaseScraper struct {
	Config *types.ScraperConfig
}

func NewBaseScraper(config *types.ScraperConfig) *BaseScraper {
	return &BaseScraper{
		Config: config,
	}
}

func (b *BaseScraper) GetName() string {
	return b.Config.Name
}

func (b *BaseScraper) GetAllowedDomains() []string {
	return b.Config.AllowedDomains
}

// CreateCollector creates a configured colly collector
func (b *BaseScraper) CreateCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(b.Config.AllowedDomains...),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited:", r.Request.URL)
	})

	return c
}
