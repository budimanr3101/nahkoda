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
	args := []string{plan.Operation, plan.Resource}

	if plan.Target != "" {
		args = append(args, plan.Target)
	}

	// Namespace
	if plan.Namespace == "all" {
		args = append(args, "-A")
	} else {
		args = append(args, "-n", plan.Namespace)
	}

	// Filter (field-selector) - SERVER SIDE
	var selectors []string
	for k, v := range plan.Filters {
		selectors = append(selectors, k+v)
	}
	if len(selectors) > 0 {
		args = append(args, "--field-selector="+strings.Join(selectors, ","))
	}

	// Execute
	cmd := exec.Command("kubectl", args...)

	// Jika tidak ada GREP, langsung stream ke stdout (native performance)
	if plan.Grep == "" {
		cmd.Stdout = os.Stdout

		// Capture stderr untuk cek "NotFound"
		var stderrBuf bytes.Buffer
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

		fmt.Printf("⚓ Menjalankan: %s\n", strings.Join(cmd.Args, " "))

		if err := cmd.Run(); err != nil {
			// Graceful error handling for NotFound
			errStr := stderrBuf.String()
			if strings.Contains(errStr, "NotFound") || strings.Contains(errStr, "not found") {
				// Resource not found is not an error in our context
				return nil
			}
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
		// Graceful error handling for NotFound
		errStr := stderrBuf.String()
		if strings.Contains(errStr, "NotFound") || strings.Contains(errStr, "not found") {
			// Resource not found is not an error in our context
			return nil
		}
		// Wrap kubectl error with context
		return errors.NewKubectlFailed(err).WithContext("command", strings.Join(cmd.Args, " "))
	}
	return nil
}
