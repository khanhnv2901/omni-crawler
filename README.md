# Omni Crawler

A flexible and extensible web scraper built in Go that supports multiple site types and output formats.

## Features

- **Modular Architecture**: Easy to add new scraper types
- **Configuration-driven**: YAML-based configuration files
- **Multiple Output Formats**: CSV support with extensible writer interface
- **Robust Scraping**: Built on top of Colly framework
- **Factory Pattern**: Clean scraper instantiation and management

## Installation

```bash
git clone https://github.com/khanhnv2901/omni-crawler.git
cd omni-crawler
go mod download
```

## Usage

### Basic Usage

```bash
# Build the scraper
go build -o scraper cmd/scraper/main.go

# List available scraper types
./scraper -list

# Run a specific scraper
./scraper -scraper ecommerce
```

### Command Line Options

- `-config-dir`: Directory containing configuration files (default: "configs")
- `-scraper`: Type of scraper to run (must match config filename)
- `-list`: List all available scraper types

## Configuration

Create YAML configuration files in the `configs/` directory. Each scraper type needs its own configuration file.

### Example Configuration (ecommerce.yaml)

```yaml
name: "ecommerce"
allowed_domains:
  - "www.scrapingcourse.com"
start_url: "https://www.scrapingcourse.com/ecommerce"
selectors:
  product_item: "li.product"
  url: "a"
  image: "img"
  name: ".product-name"
  price: ".price"
output_format: "csv"
output_file: "products.csv"
```

## Project Structure

```
omni-crawler/
├── cmd/scraper/          # Main application entry point
├── internal/
│   ├── config/           # Configuration loading
│   ├── output/           # Output writers (CSV, etc.)
│   ├── scraper/          # Base scraper and factory
│   ├── sites/            # Site-specific scrapers
│   │   └── ecommerce/    # E-commerce scraper implementation
│   └── types/            # Shared types and interfaces
└── configs/              # Configuration files
```

## Extending the Crawler

### Adding New Scraper Types

The crawler is designed to be easily extensible for different types of websites. Here's how to add support for new page types:

#### 1. Create a New Scraper Package

```bash
mkdir internal/sites/yourtype
```

#### 2. Define Your Data Structure

```go
// internal/sites/yourtype/types.go
package yourtype

type YourDataType struct {
    Field1 string `json:"field1"`
    Field2 string `json:"field2"`
    // Add fields specific to your data type
}
```

#### 3. Implement the Scraper Interface

```go
// internal/sites/yourtype/scraper.go
package yourtype

import (
    "github.com/gocolly/colly"
    "github.com/khanhnv2901/omni-crawler/internal/scraper"
    "github.com/khanhnv2901/omni-crawler/internal/types"
)

type YourScraper struct {
    *scraper.BaseScraper
    data []YourDataType
}

func NewScraper(config *types.ScraperConfig) scraper.Scraper {
    return &YourScraper{
        BaseScraper: scraper.NewBaseScraper(config),
        data:        make([]YourDataType, 0),
    }
}

func (s *YourScraper) ConfigureCollector(c *colly.Collector) error {
    c.OnHTML(s.Config.Selectors["item"], func(e *colly.HTMLElement) {
        item := YourDataType{
            Field1: e.ChildText(s.Config.Selectors["field1"]),
            Field2: e.ChildText(s.Config.Selectors["field2"]),
        }
        s.data = append(s.data, item)
    })
    return nil
}

func (s *YourScraper) Scrape(url string) ([]scraper.ScrapedData, error) {
    s.data = make([]YourDataType, 0)

    c := s.CreateCollector()
    if err := s.ConfigureCollector(c); err != nil {
        return nil, err
    }

    if err := c.Visit(url); err != nil {
        return nil, err
    }

    // Convert to interface slice
    result := make([]scraper.ScrapedData, len(s.data))
    for i, item := range s.data {
        result[i] = item
    }

    return result, nil
}
```

#### 4. Register Your Scraper

Add your scraper to `cmd/scraper/main.go`:

```go
import "github.com/khanhnv2901/omni-crawler/internal/sites/yourtype"

// In main function
factory.RegisterScraper("yourtype", yourtype.NewScraper)
```

#### 5. Create Configuration File

Create `configs/yourtype.yaml`:

```yaml
name: "yourtype"
allowed_domains:
  - "example.com"
start_url: "https://example.com/page"
selectors:
  item: ".item-selector"
  field1: ".field1-selector"
  field2: ".field2-selector"
output_format: "csv"
output_file: "yourdata.csv"
```

### Adding New Output Formats

The output system is also extensible. You can add support for JSON, XML, databases, etc.

#### 1. Implement the Writer Interface

```go
// internal/output/json.go
package output

import (
    "encoding/json"
    "os"
    "github.com/khanhnv2901/omni-crawler/internal/scraper"
)

type JSONWriter struct{}

func NewJSONWriter() Writer {
    return &JSONWriter{}
}

func (w *JSONWriter) Write(data []scraper.ScrapedData, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(data)
}
```

#### 2. Register the Writer

Add your writer to the output format switch in `cmd/scraper/main.go`:

```go
switch config.OutputFormat {
case "csv":
    writer = output.NewCSVWriter()
case "json":
    writer = output.NewJSONWriter()
case "xml":
    writer = output.NewXMLWriter()
default:
    writer = output.NewCSVWriter()
}
```

### Advanced Customization

#### Custom Base Scrapers

For complex scraping needs, you can create specialized base scrapers:

```go
// internal/scraper/advanced_base.go
type AdvancedBaseScraper struct {
    *BaseScraper
    // Add common functionality for advanced scrapers
}

func (a *AdvancedBaseScraper) CreateCollectorWithProxy() *colly.Collector {
    c := a.CreateCollector()
    // Add proxy configuration
    // Add rate limiting
    // Add custom headers
    return c
}
```

#### Database Output

For database storage, implement a database writer:

```go
type DatabaseWriter struct {
    db *sql.DB
}

func (w *DatabaseWriter) Write(data []scraper.ScrapedData, tableName string) error {
    // Insert data into database
    return nil
}
```

#### Configuration Extensions

Add custom configuration fields in `internal/types/config.go`:

```go
type ScraperConfig struct {
    Name           string            `yaml:"name"`
    AllowedDomains []string          `yaml:"allowed_domains"`
    StartURL       string            `yaml:"start_url"`
    Selectors      map[string]string `yaml:"selectors"`
    OutputFormat   string            `yaml:"output_format"`
    OutputFile     string            `yaml:"output_file"`

    // Custom fields
    RateLimit      int               `yaml:"rate_limit"`
    UserAgent      string            `yaml:"user_agent"`
    Headers        map[string]string `yaml:"headers"`
    Pagination     PaginationConfig  `yaml:"pagination"`
}
```

## Dependencies

- **Colly**: Web scraping framework
- **Go-Query**: HTML parsing
- **YAML v2**: Configuration parsing

## License

This project is open source. Feel free to use and modify as needed.