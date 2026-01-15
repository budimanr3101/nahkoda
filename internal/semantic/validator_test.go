package semantic

import (
	"testing"
)

func TestValidateResourceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid Name", "my-pod", false},
		{"Valid Number", "123", false},
		{"Valid Dot", "app.v1", false},
		{"Invalid Space", "my pod", true},
		{"Invalid Shell Injection ;", "pod;rm -rf /", true},
		{"Invalid Shell Injection |", "pod|cat /etc/passwd", true},
		{"Invalid Quote", "pod'name", true},
		{"Invalid Double Quote", "pod\"name", true},
		{"Invalid Backtick", "pod`date`", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateResourceName(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("ValidateResourceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
