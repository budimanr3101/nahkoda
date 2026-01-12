package semantic

import (
	"nahkoda/internal/parser"
	"testing"
)

func TestResolveV100Objects(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"liat armada", "deployment"},
		{"liat penjaga", "daemonset"},
		{"liat pelabuhan", "service"},
		{"liat mercusuar", "ingress"},
		{"liat peta", "configmap"},
		{"liat sandi", "secret"},
	}

	for _, tt := range tests {
		ast, _ := parser.Parse(tt.input)
		intent, err := Resolve(ast)
		if err != nil {
			t.Fatalf("Resolve failed for %s: %v", tt.input, err)
		}
		if intent.Objek != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, intent.Objek)
		}
	}
}

func TestBikinV100Objects(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"bikin armada app", "deployment"},
		{"bikin pelabuhan app", "service"},
		{"bikin sandi api-key", "secret"},
	}

	for _, tt := range tests {
		ast, _ := parser.Parse(tt.input)
		intent, err := Resolve(ast)
		if err != nil {
			t.Fatalf("Resolve failed for %s: %v", tt.input, err)
		}
		if intent.Objek != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, intent.Objek)
		}
	}
}
