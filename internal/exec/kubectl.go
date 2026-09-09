package exec

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// KubectlClient defines the interface for interacting with kubectl.
// This allows mocking for tests.
type KubectlClient interface {
	Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error
	Start(args []string, stdout, stderr io.Writer) (*exec.Cmd, error)
}

// StandardKubectlClient is the default implementation that calls the kubectl binary.
type StandardKubectlClient struct {
	KubectlPath string
	Timeout     time.Duration
}

func (c *StandardKubectlClient) getCmd() string {
	if c.KubectlPath == "" {
		return "kubectl"
	}
	return c.KubectlPath
}

func (c *StandardKubectlClient) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	// Interactive exec dan follow logs sengaja tidak dibatasi. Keduanya dirancang
	// hidup sampai pengguna keluar, sehingga timeout global justru merusak fitur.
	if c.Timeout <= 0 || isLongRunningCommand(args) {
		cmd := exec.Command(c.getCmd(), args...)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		return cmd.Run()
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.getCmd(), args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("kubectl melewati timeout %s: %w", c.Timeout, context.DeadlineExceeded)
	}
	return err
}

func isLongRunningCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	if args[0] == "exec" {
		return true
	}
	if args[0] != "logs" {
		return false
	}
	for _, arg := range args[1:] {
		if arg == "-f" || arg == "--follow" {
			return true
		}
	}
	return false
}

func (c *StandardKubectlClient) Start(args []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	cmd := exec.Command(c.getCmd(), args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd, cmd.Start()
}

// Global client instance is removed in favor of Dependency Injection
