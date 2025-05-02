package main

import (
	"fmt"
	"sort"

	"github.com/theory/jsonpath"
)

type Processor struct {
	verbose bool
}

func NewProcessor(verbose bool) *Processor {
	return &Processor{
		verbose: verbose,
	}
}

// Extracts values from data using JSONPath
func (p *Processor) ExtractValues(data interface{}, path string) ([]string, error) {
	if p.verbose {
		fmt.Printf("JSONPath extraction: %s\n", path)
	}

	// Parse the JSONPath expression
	jp, err := jsonpath.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("JSONPath parse error: %w", err)
	}

	// Execute JSONPath - returns NodeList
	result := jp.Select(data)

	// Convert NodeList to slice to work with our existing logic
	resultInterface := p.convertNodeListToInterface(result)

	// Convert results to string slice
	values := p.convertToStringSlice(resultInterface)

	if p.verbose {
		p.printExtractedValues(values)
	}

	// Sort values for consistent comparison
	sort.Strings(values)

	return values, nil
}

// Converts NodeList to interface slice
func (p *Processor) convertNodeListToInterface(result jsonpath.NodeList) []interface{} {
	var resultInterface []interface{}
	for _, r := range result {
		resultInterface = append(resultInterface, r)
	}
	return resultInterface
}

// Prints extracted values in verbose mode
func (p *Processor) printExtractedValues(values []string) {
	fmt.Printf("Number of extracted values: %d\n", len(values))
	fmt.Println("Extracted values:")
	for i, value := range values {
		fmt.Printf("  [%3d] %s\n", i, value)
	}
	fmt.Println() // Empty line for readability
}

// Converts results of various types to string slice
func (p *Processor) convertToStringSlice(result interface{}) []string {
	var values []string

	switch v := result.(type) {
	case []interface{}:
		values = p.processArray(v)
	case map[string]interface{}:
		values = p.processMap(v)
	case string:
		values = []string{v}
	case nil:
		return []string{}
	default:
		values = []string{fmt.Sprintf("%v", v)}
	}

	return values
}

// Processes array type result
func (p *Processor) processArray(arr []interface{}) []string {
	var values []string

	for _, item := range arr {
		itemValues := p.processArrayItem(item)
		values = append(values, itemValues...)
	}

	return values
}

// Processes a single array item
func (p *Processor) processArrayItem(item interface{}) []string {
	switch v := item.(type) {
	case []interface{}:
		return p.convertToStringSlice(v)
	case map[string]interface{}:
		return p.processMap(v)
	default:
		return []string{fmt.Sprintf("%v", v)}
	}
}

// Processes map type result
func (p *Processor) processMap(m map[string]interface{}) []string {
	var values []string
	for _, value := range m {
		values = append(values, fmt.Sprintf("%v", value))
	}
	return values
}
