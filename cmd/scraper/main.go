package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/khanhnv2901/omni-crawler/internal/config"
	"github.com/khanhnv2901/omni-crawler/internal/output"
	"github.com/khanhnv2901/omni-crawler/internal/scraper"
	"github.com/khanhnv2901/omni-crawler/internal/sites/ecommerce"
)

func main() {
	// Command line flags
	var (
		configDir   = flag.String("config-dir", "configs", "Directory with configs")
		scraperType = flag.String("scraper", "", "Type of scraper torun")
		listTypes   = flag.Bool("list", false, "List available scraper types")
	)
	flag.Parse()

	// List available  scrapers
	if *listTypes {
		configs, err := config.LoadAllConfigs(*configDir)
		if err != nil {
			log.Fatalf("Failed to load configurations: %v", err)
		}

		fmt.Println("Available scraper types:")
		for name := range configs {
			fmt.Printf(" - %s\n", name)
		}
		return
	}

	// Load configuration
	configPath := fmt.Sprintf("%s/%s.yaml", *configDir, *scraperType)
	config, err := config.LoadScraperConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load scraper config: %v", err)
	}

	// Create scraper factory and register scrapers
	factory := scraper.NewFactory()
	factory.RegisterScraper("ecommerce", ecommerce.NewScraper)

	// Create scraper instance
	scraperInstance, err := factory.CreateScraper(*scraperType, config)
	if err != nil {
		log.Fatalf("Failed to create scraper: %v", err)
	}

	// Scrape data
	fmt.Printf("Starting %s scraper...\n", scraperInstance.GetName())
	data, err := scraperInstance.Scrape(config.StartURL)
	if err != nil {
		log.Fatalf("Failed to scrape data: %v", err)
	}

	fmt.Printf("Scraped %d items\n", len(data))

	// Write output
	if len(data) > 0 {
		var writer output.Writer
		switch config.OutputFormat {
		case "csv":
			writer = output.NewCSVWriter()
		default:
			writer = output.NewCSVWriter()
		}

		if err := writer.Write(data, config.OutputFile); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}

		fmt.Printf("Data saved to %s\n", config.OutputFile)
	}
}
