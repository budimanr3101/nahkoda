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
	args := []string{}

	if plan.Operation == "exec" {
		args = append(args, "exec", "-it", plan.Target)
	} else if plan.Operation == "logs" {
		args = append(args, "logs", plan.Target)
	} else if plan.Operation == "run" {
		args = append(args, "run", plan.Target)
	} else {
		args = append(args, plan.Operation, plan.Resource)
		if plan.Target != "" {
			args = append(args, plan.Target)
		}
	}

	// Namespace logic
	if plan.Operation != "config" {
		if plan.Namespace == "all" {
			args = append(args, "-A")
		} else {
			args = append(args, "-n", plan.Namespace)
		}
	}

	var selectors []string
	for k, v := range plan.Filters {
		selectors = append(selectors, k+v)
	}
	if len(selectors) > 0 {
		args = append(args, "--field-selector="+strings.Join(selectors, ","))
	}

	if plan.Operation == "exec" {
		args = append(args, "--", "/bin/sh")
	}

	// Append generic Flags
	args = append(args, plan.Flags...)

	cmd := exec.Command("kubectl", args...)

	// Capture stderr for error analysis
	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if plan.Operation == "exec" {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		fmt.Printf("⚓ Menjalankan (Interactive): %s\n", strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			checkAndPrintHint(stderrBuf.String(), plan)
			return errors.NewKubectlFailed(err).WithContext("command", strings.Join(cmd.Args, " "))
		}
		return nil
	}

	if plan.Grep == "" {
		cmd.Stdout = os.Stdout
		fmt.Printf("⚓ Menjalankan: %s\n", strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			checkAndPrintHint(stderrBuf.String(), plan)
			return errors.NewKubectlFailed(err).WithContext("command", strings.Join(cmd.Args, " "))
		}
		return nil
	}

	// Client-side filtering (Grep)
	fmt.Printf("⚓ Menjalankan: %s | grep '%s' (invert=%v)\n", strings.Join(cmd.Args, " "), plan.Grep, plan.GrepInvert)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

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
