package scraper

import (
	"fmt"

	"github.com/khanhnv2901/omni-crawler/internal/types"
)

// Factory creates scrapers based on type
type Factory struct {
	scrapers map[string]func(*types.ScraperConfig) Scraper
}

// NewFactory creates a new scraper factory
func NewFactory() *Factory {
	return &Factory{
		scrapers: make(map[string]func(*types.ScraperConfig) Scraper),
	}
}

// RegisterScraper registers a scraper creator function
func (f *Factory) CreateScraper(scraperType string, config *types.ScraperConfig) (Scraper, error) {
	creator, exists := f.scrapers[scraperType]
	if !exists {
		return nil, fmt.Errorf("scraper type '%s' not found", scraperType)
	}
	return creator(config), nil
}

func (f *Factory) RegisterScraper(scraperType string, creator func(*types.ScraperConfig) Scraper) {
	f.scrapers[scraperType] = creator
}

// GetAvailableScrapers returns list of available scraper types
func (f *Factory) GetAvailableScrapers() []string {
	var types []string
	for t := range f.scrapers {
		types = append(types, t)
	}
	return types
}
