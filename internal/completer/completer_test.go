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
		{
			name:          "After 'di' must suggest 'geladak'",
			textBefore:    "liat kru di ",
			wordBefore:    "",
			expectedCount: 1,
			shouldContain: "geladak",
		},
		{
			name:          "Adaptive: 'baca' suggests 'jurnal'",
			textBefore:    "baca ",
			wordBefore:    "",
			expectedCount: 1,
			shouldContain: "jurnal",
		},
		{
			name:          "Adaptive: 'atur' suggests 'armada'",
			textBefore:    "atur ",
			wordBefore:    "",
			expectedCount: 1,
			shouldContain: "armada",
		},
		{
			name:          "Adaptive: 'masuk' suggests 'kru'",
			textBefore:    "masuk ",
			wordBefore:    "",
			expectedCount: 1,
			shouldContain: "kru",
		},
		{
			name:       "After 'geladak' suggests namespaces",
			textBefore: "liat kru di geladak ",
			wordBefore: "",
			// We don't check for 'default' because kubectl might not return anything in CI
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
