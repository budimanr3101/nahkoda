package main

import (
	"io"
	"os/exec"
	"strings"
	"testing"

	nahkodaexec "nahkoda/internal/exec"
)

type noOpKubectlClient struct{}

func (noOpKubectlClient) Run([]string, io.Reader, io.Writer, io.Writer) error {
	return nil
}

func (noOpKubectlClient) Start([]string, io.Writer, io.Writer) (*exec.Cmd, error) {
	return nil, nil
}

func TestProcessCommand_NonInteractiveSuggestion(t *testing.T) {
	executor := nahkodaexec.NewExecutor(noOpKubectlClient{})
	err := processCommand("liat kur", executor, "default", false)
	if err == nil {
		t.Fatal("processCommand() should reject typo")
	}
	if !strings.Contains(err.Error(), "mungkin maksud Kapten: kru") {
		t.Fatalf("error = %q, want non-interactive suggestion", err.Error())
	}
}
