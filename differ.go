package main

import (
	"fmt"
)

type Differ struct {
	config   *Config
	client   *Client
	jsonpath *Processor
	verbose  bool
}

type ComparisonResult struct {
	Endpoint1    string
	Endpoint2    string
	JSONPath1    string
	JSONPath2    string
	MatchedItems []string
	OnlyIn1Items []string
	OnlyIn2Items []string
	IsMatch      bool
}

func NewDiffer(cfg *Config, verbose bool) *Differ {
	return &Differ{
		config:   cfg,
		client:   NewClient(cfg.GetTimeout(), verbose),
		jsonpath: NewProcessor(verbose),
		verbose:  verbose,
	}
}

// Executes the comparison between endpoints
func (d *Differ) RunComparison() bool {
	result, err := d.compareEndpoints()
	if err != nil {
		fmt.Printf("Error during comparison: %v\n", err)
		return false
	}

	d.printResult(result)

	return result.IsMatch
}

// Compares the two endpoints
func (d *Differ) compareEndpoints() (*ComparisonResult, error) {
	result := &ComparisonResult{
		Endpoint1: d.config.Endpoint1.URL,
		Endpoint2: d.config.Endpoint2.URL,
		JSONPath1: d.config.Endpoint1.JSONPath,
		JSONPath2: d.config.Endpoint2.JSONPath,
	}

	// Merge headers
	headers1 := mergeHeaders(d.config.General.Headers, d.config.Endpoint1.Headers)
	headers2 := mergeHeaders(d.config.General.Headers, d.config.Endpoint2.Headers)

	// Fetch data from first endpoint
	data1, err := d.client.FetchJSON(d.config.Endpoint1.URL, headers1)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from endpoint1: %w", err)
	}

	// Fetch data from second endpoint
	data2, err := d.client.FetchJSON(d.config.Endpoint2.URL, headers2)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from endpoint2: %w", err)
	}

	// Extract values using JSONPath
	values1, err := d.jsonpath.ExtractValues(data1, d.config.Endpoint1.JSONPath)
	if err != nil {
		return nil, fmt.Errorf("error extracting values from endpoint1: %w", err)
	}

	values2, err := d.jsonpath.ExtractValues(data2, d.config.Endpoint2.JSONPath)
	if err != nil {
		return nil, fmt.Errorf("error extracting values from endpoint2: %w", err)
	}

	// Compare values
	result.MatchedItems, result.OnlyIn1Items, result.OnlyIn2Items = compareValues(values1, values2)
	result.IsMatch = len(result.OnlyIn1Items) == 0 && len(result.OnlyIn2Items) == 0

	return result, nil
}

// Merges two header maps
func mergeHeaders(general, specific map[string]string) map[string]string {
	merged := make(map[string]string)

	// Copy general headers
	for k, v := range general {
		merged[k] = v
	}

	// Override with specific headers
	for k, v := range specific {
		merged[k] = v
	}

	return merged
}

// Compares two string slices
func compareValues(values1, values2 []string) (matched, onlyIn1, onlyIn2 []string) {
	// Create maps of values
	set1 := buildValueSet(values1)
	set2 := buildValueSet(values2)

	// Extract matched and unique values
	matched, onlyIn1 = extractComparisonResults(set1, set2)
	_, onlyIn2 = extractComparisonResults(set2, set1)

	return matched, onlyIn1, onlyIn2
}

// Builds a set of values from a slice
func buildValueSet(values []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range values {
		set[v] = true
	}
	return set
}

// Extracts matched and unique values from two sets
func extractComparisonResults(set1, set2 map[string]bool) (common, unique []string) {
	for v := range set1 {
		if set2[v] {
			common = append(common, v)
		} else {
			unique = append(unique, v)
		}
	}
	return common, unique
}

// Displays the comparison result
func (d *Differ) printResult(result *ComparisonResult) {
	d.printHeader(result)
	d.printSummary(result)

	if d.verbose {
		d.printDetailedComparison(result)
	}

	d.printUnmatchedItems(result)
	d.printFinalStatus(result)
}

// Prints comparison header
func (d *Differ) printHeader(result *ComparisonResult) {
	fmt.Printf("\n=== REST API Comparison ===\n")
	fmt.Printf("Endpoint1: %s (%s)\n", result.Endpoint1, result.JSONPath1)
	fmt.Printf("Endpoint2: %s (%s)\n", result.Endpoint2, result.JSONPath2)
	fmt.Println("\nResults:")
}

// Prints summary of matched items
func (d *Differ) printSummary(result *ComparisonResult) {
	fmt.Printf("  Matched: %d items\n", len(result.MatchedItems))
}

// Prints detailed comparison in verbose mode
func (d *Differ) printDetailedComparison(result *ComparisonResult) {
	fmt.Println("\nDetailed comparison:")

	d.printMatchedValues(result)
	d.printComparisonDetails(result)
}

// Prints matched values in verbose mode
func (d *Differ) printMatchedValues(result *ComparisonResult) {
	if len(result.MatchedItems) == 0 {
		return
	}

	fmt.Printf("  Matched Values (both endpoints):\n")
	for i, item := range result.MatchedItems {
		fmt.Printf("    [%3d] MATCH: '%s'\n", i, item)
	}
}

// Prints comparison details for unmatched items
func (d *Differ) printComparisonDetails(result *ComparisonResult) {
	if len(result.OnlyIn1Items) == 0 && len(result.OnlyIn2Items) == 0 {
		return
	}

	fmt.Println("\n  Comparison Details:")
	maxLen := max(len(result.OnlyIn1Items), len(result.OnlyIn2Items))

	for i := 0; i < maxLen; i++ {
		compareIndex := i + len(result.MatchedItems)
		d.printComparisonRow(result, i, compareIndex)
	}
}

// Prints a single comparison row
func (d *Differ) printComparisonRow(result *ComparisonResult, idx, printIdx int) {
	left, right := "(none)", "(none)"

	if idx < len(result.OnlyIn1Items) {
		left = result.OnlyIn1Items[idx]
	}

	if idx < len(result.OnlyIn2Items) {
		right = result.OnlyIn2Items[idx]
	}

	status := getComparisonStatus(left, right)
	fmt.Printf("    [%3d] %-25s <-> %-25s | %s\n", printIdx, left, right, status)
}

// Gets comparison status based on left and right values
func getComparisonStatus(left, right string) string {
	switch {
	case left == "(none)":
		return "ONLY IN ENDPOINT2"
	case right == "(none)":
		return "ONLY IN ENDPOINT1"
	default:
		return "MISMATCH"
	}
}

// Prints unmatched items in non-verbose mode
func (d *Differ) printUnmatchedItems(result *ComparisonResult) {
	if !d.verbose {
		d.printOnlyInItems("Only in Endpoint1", result.OnlyIn1Items)
		d.printOnlyInItems("Only in Endpoint2", result.OnlyIn2Items)
	} else {
		// Just show counts in verbose mode
		if len(result.OnlyIn1Items) > 0 {
			fmt.Printf("  Only in Endpoint1: %d items\n", len(result.OnlyIn1Items))
		}
		if len(result.OnlyIn2Items) > 0 {
			fmt.Printf("  Only in Endpoint2: %d items\n", len(result.OnlyIn2Items))
		}
	}
}

// Prints items only present in one endpoint
func (d *Differ) printOnlyInItems(label string, items []string) {
	if len(items) == 0 {
		return
	}

	fmt.Printf("  %s: %d items\n", label, len(items))
	for _, item := range items {
		fmt.Printf("    - %s\n", item)
	}
}

// Prints final comparison status
func (d *Differ) printFinalStatus(result *ComparisonResult) {
	status := "MATCH"
	if !result.IsMatch {
		status = "MISMATCH"
	}
	fmt.Printf("\nComparison Status: %s\n", status)
}

// Helper function for max of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
