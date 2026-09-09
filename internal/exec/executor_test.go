package exec

import (
	stderrors "errors"
	"io"
	"nahkoda/internal/errors"
	"nahkoda/internal/planner"
	"os/exec"
	"strings"
	"testing"
	"time"
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

	// Execute should not error
	if err := executor.Execute(plan); err != nil {
		t.Errorf("Execute() error = %v", err)
	}

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

	// Execute should not error
	if err := executor.Execute(plan); err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Note: Filter order is map iteration based, might need robust check.
	// But here single key.
	if !strings.Contains(strings.Join(mock.LastArgs, " "), "--field-selector=status.phase=Running") {
		t.Errorf("Expected field selector, got %s", mock.LastArgs)
	}
}

type failingKubectlClient struct{}

func (failingKubectlClient) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if stderr != nil {
		_, _ = stderr.Write([]byte("cluster unavailable"))
	}
	return stderrors.New("cluster unavailable")
}

func (failingKubectlClient) Start(args []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	return nil, stderrors.New("cluster unavailable")
}

func TestRunAudit_PropagatesCriticalErrors(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	executor := NewExecutor(failingKubectlClient{})
	if err := executor.Execute(planner.Plan{Operation: "audit"}); err == nil {
		t.Fatal("audit should fail when core kubectl checks fail")
	}
}

type contextAwareKubectlClient struct {
	Calls     [][]string
	NoContext bool
}

func (c *contextAwareKubectlClient) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	copiedArgs := append([]string(nil), args...)
	c.Calls = append(c.Calls, copiedArgs)

	if strings.Join(args, " ") == "config current-context" {
		if c.NoContext {
			if stderr != nil {
				_, _ = stderr.Write([]byte("current-context is not set"))
			}
			return stderrors.New("exit status 1")
		}
		if stdout != nil {
			_, _ = stdout.Write([]byte("test-context\n"))
		}
	}
	return nil
}

func (c *contextAwareKubectlClient) Start(args []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	return nil, nil
}

func TestExecuteFailsFastWithoutActiveContext(t *testing.T) {
	client := &contextAwareKubectlClient{NoContext: true}
	executor := NewExecutor(client)

	err := executor.Execute(planner.Plan{
		Operation: "get",
		Resource:  "pods",
		Namespace: "all",
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want no-active-context error")
	}

	var nahkodaErr *errors.NahkodaError
	if !stderrors.As(err, &nahkodaErr) || !nahkodaErr.IsType(errors.ErrNoActiveContext) {
		t.Fatalf("Execute() error = %v, want ErrNoActiveContext", err)
	}
	if len(client.Calls) != 1 || strings.Join(client.Calls[0], " ") != "config current-context" {
		t.Fatalf("kubectl calls = %v, want only context preflight", client.Calls)
	}
}

func TestExecuteConfigBypassesContextCheck(t *testing.T) {
	client := &contextAwareKubectlClient{NoContext: true}
	executor := NewExecutor(client)

	err := executor.Execute(planner.Plan{
		Operation: "config",
		Resource:  "get-contexts",
	})
	if err != nil {
		t.Fatalf("Execute(config) error = %v", err)
	}
	if len(client.Calls) != 1 || strings.Join(client.Calls[0], " ") != "config get-contexts" {
		t.Fatalf("kubectl calls = %v, want config get-contexts without preflight", client.Calls)
	}
}

func TestExecuteDryRunBypassesContextCheck(t *testing.T) {
	client := &contextAwareKubectlClient{NoContext: true}
	executor := NewExecutor(client)
	executor.DryRun = true

	err := executor.Execute(planner.Plan{
		Operation: "get",
		Resource:  "pods",
		Namespace: "all",
	})
	if err != nil {
		t.Fatalf("Execute(dry-run) error = %v", err)
	}
	if len(client.Calls) != 0 {
		t.Fatalf("dry-run kubectl calls = %v, want none", client.Calls)
	}
}

type delayedKubectlClient struct {
	Delay time.Duration
}

func (c delayedKubectlClient) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if strings.Join(args, " ") == "config current-context" {
		if stdout != nil {
			_, _ = stdout.Write([]byte("test-context\n"))
		}
		return nil
	}
	time.Sleep(c.Delay)
	return nil
}

func (c delayedKubectlClient) Start(args []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	return nil, nil
}

func TestExecuteReportsSlowFiniteCommand(t *testing.T) {
	var notice strings.Builder
	executor := NewExecutor(delayedKubectlClient{Delay: 20 * time.Millisecond})
	executor.slowNoticeDelay = time.Millisecond
	executor.noticeWriter = &notice

	err := executor.Execute(planner.Plan{
		Operation: "get",
		Resource:  "pods",
		Namespace: "all",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(notice.String(), "belum merespons") {
		t.Fatalf("slow-command notice = %q, want progress feedback", notice.String())
	}
}

func TestExecuteDoesNotReportSlowNoticeForFollowedLogs(t *testing.T) {
	var notice strings.Builder
	executor := NewExecutor(delayedKubectlClient{Delay: 20 * time.Millisecond})
	executor.slowNoticeDelay = time.Millisecond
	executor.noticeWriter = &notice

	err := executor.Execute(planner.Plan{
		Operation: "logs",
		Target:    "api-7",
		Namespace: "default",
		Flags:     []string{"-f"},
	})
	if err != nil {
		t.Fatalf("Execute(logs -f) error = %v", err)
	}
	if notice.Len() != 0 {
		t.Fatalf("logs -f notice = %q, want no finite-command warning", notice.String())
	}
}
