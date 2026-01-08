package exec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"nahkoda/internal/errors"
	"nahkoda/internal/planner"
)

func Execute(plan planner.Plan) error {
	// Construct args: kubectl [operation] [resource] [target]
	// ARGS CONSTRUCTION
	args := []string{}

	if plan.Operation == "exec" {
		// exec specific: kubectl exec -it [target] [-n namespace] -- /bin/sh
		args = append(args, "exec", "-it", plan.Target)
	} else if plan.Operation == "logs" {
		// logs specific: kubectl logs [target] [-n namespace]
		args = append(args, "logs", plan.Target)
	} else {
		// default: kubectl [operation] [resource] [target]
		args = append(args, plan.Operation, plan.Resource)
		if plan.Target != "" {
			args = append(args, plan.Target)
		}
	}

	// Namespace (SKIP if operation is "config")
	// "config" commands (get-contexts, use-context) do not support -n/-A
	if plan.Operation != "config" {
		if plan.Namespace == "all" {
			args = append(args, "-A")
		} else {
			args = append(args, "-n", plan.Namespace)
		}
	}

	// Filter (field-selector) - SERVER SIDE
	var selectors []string
	for k, v := range plan.Filters {
		selectors = append(selectors, k+v)
	}
	if len(selectors) > 0 {
		args = append(args, "--field-selector="+strings.Join(selectors, ","))
	}

	// Finalize Exec Args
	if plan.Operation == "exec" {
		args = append(args, "--", "/bin/sh")
	}

	// Append generic Flags
	args = append(args, plan.Flags...)

	// Execute
	cmd := exec.Command("kubectl", args...)

	// INTERACTIVE MODE (EXEC)
	if plan.Operation == "exec" {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("⚓ Menjalankan (Interactive): %s\n", strings.Join(cmd.Args, " "))
		return cmd.Run()
	}

	// Jika tidak ada GREP, langsung stream ke stdout (native performance)
	if plan.Grep == "" {
		cmd.Stdout = os.Stdout

		// Capture stderr untuk cek "NotFound"
		var stderrBuf bytes.Buffer
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

		fmt.Printf("⚓ Menjalankan: %s\n", strings.Join(cmd.Args, " "))

		if err := cmd.Run(); err != nil {
			// Wrap kubectl error with context
			return errors.NewKubectlFailed(err).WithContext("command", strings.Join(cmd.Args, " "))
		}
		return nil
	}

	// JIKA ADA GREP → Capture & Filter (Client Side)
	fmt.Printf("⚓ Menjalankan: %s | grep '%s' (invert=%v)\n", strings.Join(cmd.Args, " "), plan.Grep, plan.GrepInvert)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Capture stderr juga disini
	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()

		// Header (baris pertama) selalu tampil
		if lineCount == 0 {
			fmt.Println(line)
			lineCount++
			continue
		}

		// Filter Logic
		var match bool
		if plan.GrepRegex {
			match, _ = regexp.MatchString(plan.Grep, line)
		} else {
			match = strings.Contains(line, plan.Grep)
		}

		if plan.GrepInvert {
			if !match {
				fmt.Println(line)
			}
		} else {
			if match {
				fmt.Println(line)
			}
		}
		lineCount++
	}

	if err := cmd.Wait(); err != nil {
		// Wrap kubectl error with context
		return errors.NewKubectlFailed(err).WithContext("command", strings.Join(cmd.Args, " "))
	}
	return nil
}
