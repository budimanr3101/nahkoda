package exec

import (
	"context"
	stderrors "errors"
	"io"
	"os"
	"testing"
	"time"
)

func TestStandardKubectlClientTimeout(t *testing.T) {
	t.Setenv("GO_WANT_NAHKODA_HELPER", "1")
	client := &StandardKubectlClient{
		KubectlPath: os.Args[0],
		Timeout:     25 * time.Millisecond,
	}

	err := client.Run([]string{"-test.run=TestKubectlHelperProcess"}, nil, io.Discard, io.Discard)
	if !stderrors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Run() error = %v, want context deadline exceeded", err)
	}
}

func TestKubectlHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_NAHKODA_HELPER") != "1" {
		return
	}
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
}

func TestIsLongRunningCommand(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		{args: []string{"exec", "-it", "pod"}, want: true},
		{args: []string{"logs", "pod", "-f"}, want: true},
		{args: []string{"logs", "pod", "--follow"}, want: true},
		{args: []string{"logs", "pod"}, want: false},
		{args: []string{"get", "pods"}, want: false},
	}
	for _, tt := range tests {
		if got := isLongRunningCommand(tt.args); got != tt.want {
			t.Errorf("isLongRunningCommand(%v) = %v, want %v", tt.args, got, tt.want)
		}
	}
}
