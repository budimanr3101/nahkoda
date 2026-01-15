package exec

import (
	"io"
	"os/exec"
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
}

func (c *StandardKubectlClient) getCmd() string {
	if c.KubectlPath == "" {
		return "kubectl"
	}
	return c.KubectlPath
}

func (c *StandardKubectlClient) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	cmd := exec.Command(c.getCmd(), args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func (c *StandardKubectlClient) Start(args []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	cmd := exec.Command(c.getCmd(), args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd, cmd.Start()
}

// Global client instance is removed in favor of Dependency Injection
