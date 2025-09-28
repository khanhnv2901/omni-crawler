package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/khanhnv2901/omni-crawler/internal/types"
	"gopkg.in/yaml.v2"
)

// Load ScraperConfig loads a single config file
func LoadScraperConfig(configPath string) (*types.ScraperConfig, error) {
	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config types.ScraperConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// LoadAllConfigs loads all config files from a directory
func LoadAllConfigs(configDir string) (map[string]*types.ScraperConfig, error) {
	configs := make(map[string]*types.ScraperConfig)

	// Find all .yaml files
	files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to find config files: %w", err)
	}

	fmt.Printf("Files found: %+v\n", files)

	// Load each file
	for _, file := range files {
		config, err := LoadScraperConfig(file)
		if err != nil {
			return nil, fmt.Errorf("failed to load config %s: %w", file, err)
		}

		// Use filename (without extension) as key
		name := filepath.Base(file)
		name = name[:len(name)-5] // remove .yaml
		configs[name] = config
	}

	return configs, nil
}
