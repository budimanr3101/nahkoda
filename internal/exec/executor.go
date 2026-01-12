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
	if plan.Operation == "audit" {
		return runAudit()
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

func runAudit() error {
	fmt.Println("🩺 Memulai Audit Kesehatan Kapal (Health Audit)...")
	fmt.Println("=================================================")

	// 1. Kru Bermasalah
	fmt.Print("📋 Memeriksa Kru (Pods)... ")
	cmdPods := exec.Command("kubectl", "get", "pods", "-A", "--field-selector", "status.phase!=Running,status.phase!=Succeeded")
	outPods, _ := cmdPods.Output()
	linesPods := strings.Split(strings.TrimSpace(string(outPods)), "\n")
	if len(linesPods) <= 1 {
		fmt.Println("✅ Semua kru sehat.")
	} else {
		fmt.Printf("⚠️  Ditemukan %d kru bermasalah:\n", len(linesPods)-1)
		fmt.Println(string(outPods))
	}

	// 2. Mesin Mogok
	fmt.Print("⚙️  Memeriksa Mesin (Nodes)... ")
	cmdNodes := exec.Command("kubectl", "get", "nodes")
	outNodes, _ := cmdNodes.Output()
	if strings.Contains(string(outNodes), "NotReady") {
		fmt.Println("⚠️  Ada mesin yang mogok (NotReady):")
		// Grep NotReady lines
		lines := strings.Split(string(outNodes), "\n")
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
	cmdEvents := exec.Command("kubectl", "get", "events", "-A", "--field-selector", "type=Warning", "--sort-by=.metadata.creationTimestamp")
	outEvents, _ := cmdEvents.Output()
	linesEvents := strings.Split(strings.TrimSpace(string(outEvents)), "\n")
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
	cmdTop := exec.Command("kubectl", "top", "nodes")
	if err := cmdTop.Run(); err != nil {
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
