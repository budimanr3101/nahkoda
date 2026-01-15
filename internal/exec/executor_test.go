package exec

import (
	"io"
	"nahkoda/internal/planner"
	"os/exec"
	"strings"
	"testing"
)

// MockKubectlClient captures commands for testing
type MockKubectlClient struct {
	LastArgs []string
}

func (m *MockKubectlClient) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	m.LastArgs = args
	// Simulate output to avoid scanner panic or empty checks
	if stdout != nil {
		stdout.Write([]byte("NAME STATUS\npod-1 Running\n"))
	}
	return nil
}

func (m *MockKubectlClient) Start(args []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	m.LastArgs = args
	return nil, nil
}

func TestExecute_GeneratesCorrectCommand(t *testing.T) {
	mock := &MockKubectlClient{}
	// Use DI instead of global replacement
	executor := NewExecutor(mock)

	plan := planner.Plan{
		Operation: "get",
		Resource:  "pods",
		Namespace: "default",
	}

	executor.Execute(plan)

	expected := "get pods -n default"
	actual := strings.Join(mock.LastArgs, " ")

	if actual != expected {
		t.Errorf("Expected command '%s', got '%s'", expected, actual)
	}
}

func TestExecute_WithFilter(t *testing.T) {
	mock := &MockKubectlClient{}
	executor := NewExecutor(mock)

	plan := planner.Plan{
		Operation: "get",
		Resource:  "pods",
		Namespace: "default",
		Filters:   map[string]string{"status.phase=": "Running"},
	}

	executor.Execute(plan)

	// Note: Filter order is map iteration based, might need robust check.
	// But here single key.
	if !strings.Contains(strings.Join(mock.LastArgs, " "), "--field-selector=status.phase=Running") {
		t.Errorf("Expected field selector, got %s", mock.LastArgs)
	}
}
