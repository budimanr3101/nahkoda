package errors

import (
	"errors"
	"testing"
)

func TestNahkodaError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *NahkodaError
		expected string
	}{
		{
			name:     "simple error",
			err:      New(ErrUnknownWord, "kata tidak dikenali: \"xyz\""),
			expected: "kata tidak dikenali: \"xyz\"",
		},
		{
			name:     "wrapped error",
			err:      Wrap(ErrKubectlFailed, "perintah kubectl gagal", errors.New("exit status 1")),
			expected: "perintah kubectl gagal: exit status 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNahkodaError_IsType(t *testing.T) {
	err := New(ErrUnknownWord, "test")

	if !err.IsType(ErrUnknownWord) {
		t.Error("IsType() should return true for matching type")
	}

	if err.IsType(ErrUnknownAction) {
		t.Error("IsType() should return false for non-matching type")
	}
}

func TestNahkodaError_WithContext(t *testing.T) {
	err := New(ErrKubectlFailed, "test").
		WithContext("command", "kubectl get pods").
		WithContext("namespace", "default")

	if err.Context["command"] != "kubectl get pods" {
		t.Error("Context not set correctly")
	}

	if err.Context["namespace"] != "default" {
		t.Error("Context not set correctly")
	}
}

func TestNewUnknownWord(t *testing.T) {
	err := NewUnknownWord("xyz")

	if err.Type != ErrUnknownWord {
		t.Error("Wrong error type")
	}

	expected := "kata tidak dikenali: \"xyz\""
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestNewMissingTarget(t *testing.T) {
	err := NewMissingTarget("mesin")

	if err.Type != ErrMissingTarget {
		t.Error("Wrong error type")
	}

	expected := "cek mesin butuh nama mesin"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestIsResourceNotFound(t *testing.T) {
	notFoundErr := New(ErrResourceNotFound, "not found")
	otherErr := New(ErrUnknownWord, "unknown")
	stdErr := errors.New("standard error")

	if !IsResourceNotFound(notFoundErr) {
		t.Error("Should detect resource not found error")
	}

	if IsResourceNotFound(otherErr) {
		t.Error("Should not detect non-resource-not-found error")
	}

	if IsResourceNotFound(stdErr) {
		t.Error("Should not detect standard error as resource not found")
	}
}

func TestNahkodaError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	wrapped := Wrap(ErrKubectlFailed, "kubectl failed", underlying)

	if errors.Unwrap(wrapped) != underlying {
		t.Error("Unwrap should return underlying error")
	}

	simple := New(ErrUnknownWord, "test")
	if errors.Unwrap(simple) != nil {
		t.Error("Unwrap should return nil for non-wrapped error")
	}
}
