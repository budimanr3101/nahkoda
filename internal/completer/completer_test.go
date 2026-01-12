package completer

import (
	"testing"
)

func TestCompleter(t *testing.T) {
	tests := []struct {
		name          string
		textBefore    string
		wordBefore    string
		expectedCount int // checking if we get suggestions
		shouldContain string
	}{
		{
			name:          "Start of line suggests actions",
			textBefore:    "",
			wordBefore:    "",
			expectedCount: len(actions),
			shouldContain: "liat",
		},
		{
			name:          "After an action suggests objects",
			textBefore:    "liat ",
			wordBefore:    "",
			expectedCount: len(objects),
			shouldContain: "kru",
		},
		{
			name:          "Typing action prefix filters",
			textBefore:    "li",
			wordBefore:    "li",
			expectedCount: 1,
			shouldContain: "liat",
		},
		{
			name:          "Typing object prefix filters",
			textBefore:    "liat kr",
			wordBefore:    "kr",
			expectedCount: 1,
			shouldContain: "kru",
		},
		{
			name:          "After object suggests keywords",
			textBefore:    "liat kru ",
			wordBefore:    "",
			shouldContain: "di",
		},
		{
			name:          "Multiple spaces handling",
			textBefore:    "liat   ",
			wordBefore:    "",
			shouldContain: "kru",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := GetSuggestions(tt.textBefore, tt.wordBefore)

			if tt.expectedCount > 0 && len(suggestions) != tt.expectedCount {
				t.Errorf("expected %d suggestions, got %d", tt.expectedCount, len(suggestions))
			}

			found := false
			for _, s := range suggestions {
				if s.Text == tt.shouldContain {
					found = true
					break
				}
			}
			if tt.shouldContain != "" && !found {
				t.Errorf("suggestions did not contain %q", tt.shouldContain)
			}
		})
	}
}
