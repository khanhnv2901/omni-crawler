package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"

	"github.com/khanhnv2901/omni-crawler/internal/scraper"
)

type CSVWriter struct{}

func NewCSVWriter() *CSVWriter {
	return &CSVWriter{}
}

func (w *CSVWriter) Write(data []scraper.ScrapedData, filename string) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to write")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close() // Close file when function exits

	writer := csv.NewWriter(file)
	defer writer.Flush() // Flush when function exits

	// Get headers from first item using reflection
	headers := getStructFields(data[0])
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data rows
	for _, item := range data {
		record := getStructValues(item)
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// getStructFields uses reflection to get field names
func getStructFields(item interface{}) []string {
	val := reflect.ValueOf(item)
	typ := reflect.TypeOf(item)

	// Handle pointers
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	var fields []string
	for i := 0; i < val.NumField(); i++ {
		fields = append(fields, typ.Field(i).Name)
	}
	return fields
}

// getStructValues uses reflection to get field values
func getStructValues(item interface{}) []string {
	val := reflect.ValueOf(item)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var values []string
	for i := 0; i < val.NumField(); i++ {
		values = append(values, fmt.Sprintf("%v", val.Field(i).Interface()))
	}
	return values
}