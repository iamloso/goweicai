package gowencai

import (
	"testing"
)

func TestGetRandomUserAgent(t *testing.T) {
	ua := GetRandomUserAgent()
	if ua == "" {
		t.Error("Expected non-empty user agent")
	}
}

func TestParseURLParams(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected int
	}{
		{"empty url", "", 0},
		{"url with params", "http://example.com?a=1&b=2", 2},
		{"url without params", "http://example.com", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseURLParams(tt.url)
			if len(result) != tt.expected {
				t.Errorf("Expected %d params, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "value",
		},
		"arr": []interface{}{
			"first",
			"second",
		},
	}

	tests := []struct {
		name     string
		path     string
		expected interface{}
	}{
		{"simple path", "a.b", "value"},
		{"array access", "arr.0", "first"},
		{"invalid path", "x.y", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetValue(data, tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestQueryOptions(t *testing.T) {
	opts := &QueryOptions{
		Query:     "test",
		Cookie:    "test_cookie",
		QueryType: "stock",
		Page:      1,
		PerPage:   100,
		Retry:     10,
	}

	if opts.Query != "test" {
		t.Error("Query not set correctly")
	}
	if opts.Cookie != "test_cookie" {
		t.Error("Cookie not set correctly")
	}
}
