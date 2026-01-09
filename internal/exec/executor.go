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
	} else if plan.Operation == "run" {
		// run specific: kubectl run [target] ... (Resource is empty)
		args = append(args, "run", plan.Target)
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

	// Execute target (interactive or native stdout)
	if plan.Operation == "exec" {
		cmd.Stdin = os.Stdin // Interactive needs Stdin
	}

	// Capture stderr for error analysis
	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	// Interactive exec uses direct IO
	if plan.Operation == "exec" {
		cmd.Stdout = os.Stdout
		fmt.Printf("⚓ Menjalankan (Interactive): %s\n", strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			checkAndPrintHint(stderrBuf.String(), plan)
			return errors.NewKubectlFailed(err).WithContext("command", strings.Join(cmd.Args, " "))
		}
		return nil
	}

	// Native execution (without grep)
	if plan.Grep == "" {
		cmd.Stdout = os.Stdout
		fmt.Printf("⚓ Menjalankan: %s\n", strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			checkAndPrintHint(stderrBuf.String(), plan)
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
	// Capture stderr juga disini (reuse stderrBuf if needed, but it's cleaner to use a new one or reset)
	// Since we are in a different block (Grep mode), we can just use a new buffer variable name
	// to avoid confusion/shadowing issues with the previous block if we merged scopes.
	// BUT, here we are in the same function scope.
	// The previous `stderrBuf` was declared in the function scope (line 72).
	// So we should just reset it or reuse it.
	stderrBuf.Reset()
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
		checkAndPrintHint(stderrBuf.String(), plan)
		return errors.NewKubectlFailed(err).WithContext("command", strings.Join(cmd.Args, " "))
	}
	return nil
}

func checkAndPrintHint(errStr string, plan planner.Plan) {
	if strings.Contains(errStr, "NotFound") || strings.Contains(errStr, "not found") {
		if plan.Namespace == "" || plan.Namespace == "default" {
			fmt.Printf("\n💡 TIPS: Resource '%s' tidak ditemukan di geladak 'default'.\n", plan.Target)
			fmt.Println("   Coba cek di geladak lain dengan menambahkan: 'di geladak [nama]'")
		}
	}
}
