package main

import (
	"reflect"
	"testing"
)

func TestProcessorConvertToStringSlice(t *testing.T) {
	processor := NewProcessor(false)

	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name: "array of strings",
			input: []interface{}{
				"title1",
				"title2",
				"title3",
			},
			expected: []string{"title1", "title2", "title3"},
		},
		{
			name: "nested array",
			input: []interface{}{
				[]interface{}{"nested1", "nested2"},
				"single",
			},
			// This test will reveal the current problem - nested arrays are not properly handled
			expected: []string{"nested1", "nested2", "single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.convertToStringSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("convertToStringSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessorExtractValues(t *testing.T) {
	processor := NewProcessor(false)

	// Test with a simple JSON structure similar to JSONPlaceholder
	data := []interface{}{
		map[string]interface{}{
			"id":    1,
			"title": "Title 1",
		},
		map[string]interface{}{
			"id":    2,
			"title": "Title 2",
		},
		map[string]interface{}{
			"id":    3,
			"title": "Title 3",
		},
	}

	tests := []struct {
		name     string
		data     interface{}
		path     string
		expected []string
	}{
		{
			name:     "extract titles from array",
			data:     data,
			path:     "$[*].title",
			expected: []string{"Title 1", "Title 2", "Title 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.ExtractValues(tt.data, tt.path)
			if err != nil {
				t.Errorf("ExtractValues() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractValues() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIssueWithJSONPath(t *testing.T) {
	// This test simulates the issue with the JSONPlaceholder API
	processor := NewProcessor(true)

	// Sample data mimicking JSONPlaceholder response
	data := []interface{}{
		map[string]interface{}{
			"userId":    1,
			"id":        1,
			"title":     "delectus aut autem",
			"completed": false,
		},
		map[string]interface{}{
			"userId":    1,
			"id":        2,
			"title":     "quis ut nam facilis et officia qui",
			"completed": false,
		},
	}

	result, err := processor.ExtractValues(data, "$[*].title")
	if err != nil {
		t.Errorf("ExtractValues() error = %v", err)
		return
	}

	t.Logf("Number of extracted values: %d", len(result))
	t.Logf("Extracted values: %v", result)

	// We expect 2 titles, not 1
	if len(result) != 2 {
		t.Errorf("Expected 2 titles, got %d", len(result))
	}
}
