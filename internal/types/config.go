package types

type ScraperConfig struct {
	Name           string            `yaml:"name"`            // Scraper name
	AllowedDomains []string          `yaml:"allowed_domains"` // Domains to scrape
	StartURL       string            `yaml:"start_url"`       // Starting URL
	Selectors      map[string]string `yaml:"selectors"`       // CSS selectors
	OutputFormat   string            `yaml:"output_format"`   // OUtput format (csv, json)
	OutputFile     string            `yaml:"output_file"`     // Output filename
}
