package exec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"nahkoda/internal/errors"
	"nahkoda/internal/logger"
	"nahkoda/internal/planner"
)

// Executor handles the execution of plans using a KubectlClient.
type Executor struct {
	client  KubectlClient
	DryRun  bool
	Verbose bool
}

// NewExecutor creates a new Executor with the given client.
func NewExecutor(client KubectlClient) *Executor {
	return &Executor{client: client}
}

func (e *Executor) Execute(plan planner.Plan) error {
	if plan.Operation == "audit" {
		return e.runAudit()
	}

	args := []string{}

	if plan.Operation == "exec" {
		args = append(args, "exec", "-it", plan.Target)
	} else if plan.Operation == "logs" {
		args = append(args, "logs", plan.Target)
	} else if plan.Operation == "run" {
		args = append(args, "run", plan.Target)
	} else {
		ops := strings.Fields(plan.Operation)
		args = append(args, ops...)
		args = append(args, plan.Resource)
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

	// VERBOSE MODE
	if e.Verbose {
		fmt.Printf("[VERBOSE] Constructed command: kubectl %s\n", strings.Join(args, " "))
		if plan.Grep != "" {
			fmt.Printf("[VERBOSE] Client-side filter: grep '%s' (invert=%v)\n", plan.Grep, plan.GrepInvert)
		}
	}

	// DRY-RUN CHECK
	if e.DryRun {
		fmt.Printf("⚓ [DRY-RUN] Akan menjalankan: kubectl %s\n", strings.Join(args, " "))
		if plan.Grep != "" {
			fmt.Printf("                          | grep '%s'\n", plan.Grep)
		}
		return nil
	}

	// Capture stderr for error analysis
	var stderrBuf bytes.Buffer

	if plan.Operation == "exec" {
		// Run interactive
		fmt.Printf("⚓ Menjalankan (Interactive): kubectl %s\n", strings.Join(args, " "))

		// Use os.Stderr/Stdout directly for interactive mode
		if err := e.client.Run(args, os.Stdin, os.Stdout, io.MultiWriter(os.Stderr, &stderrBuf)); err != nil {
			logger.LogError(err, map[string]interface{}{"command": "kubectl " + strings.Join(args, " "), "type": "interactive"})
			checkAndPrintHint(stderrBuf.String(), plan)
			return errors.NewKubectlFailed(err).WithContext("command", "kubectl "+strings.Join(args, " "))
		}
		return nil
	}

	if plan.Grep == "" {
		fmt.Printf("⚓ Menjalankan: kubectl %s\n", strings.Join(args, " "))
		if err := e.client.Run(args, nil, os.Stdout, io.MultiWriter(os.Stderr, &stderrBuf)); err != nil {
			logger.LogError(err, map[string]interface{}{"command": "kubectl " + strings.Join(args, " ")})
			checkAndPrintHint(stderrBuf.String(), plan)
			return errors.NewKubectlFailed(err).WithContext("command", "kubectl "+strings.Join(args, " "))
		}
		return nil
	}

	// Client-side filtering (Grep)
	fmt.Printf("⚓ Menjalankan: kubectl %s | grep '%s' (invert=%v)\n", strings.Join(args, " "), plan.Grep, plan.GrepInvert)

	// Better approach for filtering: Read all stdout into buffer, then scan it.
	var stdoutBuf bytes.Buffer
	var sharedStderr bytes.Buffer

	err := e.client.Run(args, nil, &stdoutBuf, &sharedStderr)

	// Even if err != nil, we might have output to grep? Usually not with kubectl.
	if err != nil {
		logger.LogError(err, map[string]interface{}{"command": "kubectl " + strings.Join(args, " "), "grep": plan.Grep})
		checkAndPrintHint(sharedStderr.String(), plan)
		return errors.NewKubectlFailed(err).WithContext("command", "kubectl "+strings.Join(args, " "))
	}

	// Process output from buffer
	scanner := bufio.NewScanner(&stdoutBuf)
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

	return nil
}

func (e *Executor) runAudit() error {
	if e.DryRun {
		fmt.Println("⚓ [DRY-RUN] Akan menjalankan audit:")
		fmt.Println("   - kubectl get pods -A --field-selector=status.phase!=Running,status.phase!=Succeeded")
		fmt.Println("   - kubectl get nodes")
		fmt.Println("   - kubectl get events -A --field-selector=type=Warning")
		fmt.Println("   - kubectl top nodes")
		return nil
	}

	fmt.Println("🩺 Memulai Audit Kesehatan Kapal (Health Audit)...")
	fmt.Println("=================================================")

	var outBuf bytes.Buffer

	// 1. Kru Bermasalah
	fmt.Print("📋 Memeriksa Kru (Pods)... ")
	outBuf.Reset()
	e.client.Run([]string{"get", "pods", "-A", "--field-selector", "status.phase!=Running,status.phase!=Succeeded"}, nil, &outBuf, nil)
	linesPods := strings.Split(strings.TrimSpace(outBuf.String()), "\n")

	if len(linesPods) <= 1 && (outBuf.Len() == 0 || linesPods[0] == "") {
		// Output empty or just header? Kubectl get usually outputs "No resources found" to stderr if none.
		// If success and empty, implies no issues?
		// Actually if no resources, kubectl output might be empty or specific message.
		fmt.Println("✅ Semua kru sehat.")
	} else if len(linesPods) <= 1 {
		// Maybe just header
		fmt.Println("✅ Semua kru sehat.")
	} else {
		fmt.Printf("⚠️  Ditemukan %d kru bermasalah:\n", len(linesPods)-1)
		fmt.Println(outBuf.String())
	}

	// 2. Mesin Mogok
	fmt.Print("⚙️  Memeriksa Mesin (Nodes)... ")
	outBuf.Reset()
	e.client.Run([]string{"get", "nodes"}, nil, &outBuf, nil)
	if strings.Contains(outBuf.String(), "NotReady") {
		fmt.Println("⚠️  Ada mesin yang mogok (NotReady):")
		// Grep NotReady lines
		lines := strings.Split(outBuf.String(), "\n")
		for _, l := range lines {
			if strings.Contains(l, "NotReady") {
				fmt.Println("   - " + l)
			}
		}
	} else {
		fmt.Println("✅ Semua mesin siap berlayar.")
	}

	// 3. Berita Buruk (Events)
	fmt.Print("📢 Memeriksa Berita Buruk (Warning Events)... ")
	outBuf.Reset()
	e.client.Run([]string{"get", "events", "-A", "--field-selector", "type=Warning", "--sort-by=.metadata.creationTimestamp"}, nil, &outBuf, nil)
	linesEvents := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
	if len(linesEvents) <= 1 {
		fmt.Println("✅ Tidak ada berita buruk baru.")
	} else {
		fmt.Printf("⚠️  Ditemukan %d peringatan terbaru:\n", len(linesEvents)-1)
		// Show last 3 events
		slice := linesEvents
		if len(slice) > 4 {
			slice = slice[len(slice)-3:]
		}
		for _, e := range slice {
			fmt.Println("   - " + e)
		}
	}

	// 4. Beban (optional)
	fmt.Print("📊 Memeriksa Beban (Metrics)... ")
	err := e.client.Run([]string{"top", "nodes"}, nil, os.Stdout, nil) // Direct to stdout
	if err != nil {
		fmt.Println("⚠️  Layanan metrics tidak tersedia.")
	}

	fmt.Println("=================================================")
	fmt.Println("⚓ Audit selesai. Tetap waspada, Kapten!")
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
