package errors

import "fmt"

// ErrorType represents the category of error
type ErrorType int

const (
	// Parser errors
	ErrUnknownWord ErrorType = iota
	ErrMissingAction
	ErrInvalidSyntax

	// Semantic errors
	ErrUnknownAction
	ErrUnknownObject
	ErrUnknownCondition
	ErrMissingTarget

	// Executor errors
	ErrKubectlFailed
	ErrKubectlNotFound
	ErrResourceNotFound
)

// NahkodaError is a structured error type
type NahkodaError struct {
	Type       ErrorType
	Message    string
	Suggestion string
	Context    map[string]string
	Err        error // underlying error
}

// Error implements the error interface
func (e *NahkodaError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *NahkodaError) Unwrap() error {
	return e.Err
}

// IsType checks if error is of specific type
func (e *NahkodaError) IsType(t ErrorType) bool {
	return e.Type == t
}

// WithContext adds context to the error
func (e *NahkodaError) WithContext(key, value string) *NahkodaError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// New creates a new NahkodaError
func New(errType ErrorType, message string) *NahkodaError {
	return &NahkodaError{
		Type:    errType,
		Message: message,
		Context: make(map[string]string),
	}
}

// Wrap wraps an existing error with NahkodaError
func Wrap(errType ErrorType, message string, err error) *NahkodaError {
	return &NahkodaError{
		Type:    errType,
		Message: message,
		Context: make(map[string]string),
		Err:     err,
	}
}

// Helper functions for common errors

func NewUnknownWord(word string) *NahkodaError {
	return New(ErrUnknownWord, fmt.Sprintf("kata tidak dikenali: %q", word))
}

func (e *NahkodaError) WithSuggestion(suggestion string) *NahkodaError {
	e.Suggestion = suggestion
	return e
}

func NewUnknownAction() *NahkodaError {
	return New(ErrUnknownAction, "aksi tidak dikenali")
}

func NewUnknownObject() *NahkodaError {
	return New(ErrUnknownObject, "objek tidak dikenali")
}

func NewUnknownCondition(condition string) *NahkodaError {
	return New(ErrUnknownCondition, fmt.Sprintf("kondisi tidak dikenali: %s", condition))
}

func NewMissingTarget(objek string) *NahkodaError {
	return New(ErrMissingTarget, fmt.Sprintf("cek %s butuh nama %s", objek, objek))
}

func NewKubectlFailed(err error) *NahkodaError {
	return Wrap(ErrKubectlFailed, "perintah kubectl gagal", err)
}

func NewResourceNotFound(resource string) *NahkodaError {
	return New(ErrResourceNotFound, fmt.Sprintf("resource %q tidak ditemukan", resource))
}

// IsResourceNotFound checks if error is a resource not found error
func IsResourceNotFound(err error) bool {
	if ne, ok := err.(*NahkodaError); ok {
		return ne.Type == ErrResourceNotFound
	}
	return false
}
