package main

import (
	"fmt"
	"os"
	"strings"

	"nahkoda/internal/exec"
	"nahkoda/internal/parser"
	"nahkoda/internal/planner"
	"nahkoda/internal/semantic"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	// Check for help flags
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		return
	}

	input := strings.Join(os.Args[1:], " ")

	// 1️⃣ PARSER
	ast, err := parser.Parse(input)
	if err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}

	// 2️⃣ SEMANTIC (STRICT)
	intent, err := semantic.Resolve(ast)
	if err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}

	// 3️⃣ PLANNER
	plan := planner.Build(intent)

	// 4️⃣ EXECUTOR
	if err := exec.Execute(plan); err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}
}

func printHelp() {
	helpText := `⚓ Nahkoda — Bahasa manusia untuk Kubernetes

Cara pakai:
  nahkoda [kalimat perintah]

Contoh:
  nahkoda liat kru
  nahkoda liat kru rusak
  nahkoda liat kru di geladak auth
  nahkoda baca jurnal healthy-pod-1
  nahkoda masuk healthy-pod-1
  nahkoda liat berita

Opsi:
  -h, --help    Tampilkan bantuan ini
`
	fmt.Println(helpText)
}
