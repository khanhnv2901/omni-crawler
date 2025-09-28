package output

import "github.com/khanhnv2901/omni-crawler/internal/scraper"

type Writer interface {
	Write(data []scraper.ScrapedData, filename string) error
}