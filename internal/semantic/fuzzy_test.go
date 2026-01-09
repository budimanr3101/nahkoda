package semantic

import "testing"

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s, t string
		want int
	}{
		{"kru", "kur", 2},
		{"liat", "lite", 2},
		{"sehat", "sehet", 1},
		{"kru", "kru", 0},
		{"", "abc", 3},
		{"abc", "", 3},
	}

	for _, tt := range tests {
		if got := LevenshteinDistance(tt.s, tt.t); got != tt.want {
			t.Errorf("LevenshteinDistance(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
		}
	}
}

func TestFindSuggestion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"kur", "kru"},
		{"liatt", "liat"},
		{"ruusak", "rusak"},
		{"sehet", "sehat"},
		{"xyz", ""},        // Terlalu beda
		{"k", ""},          // Terlalu pendek
		{"kapal", "kapal"}, // Sudah pas (tapi fungsi ini biasanya dipanggil jika unknown)
	}

	for _, tt := range tests {
		if got := FindSuggestion(tt.input); got != tt.want {
			t.Errorf("FindSuggestion(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
