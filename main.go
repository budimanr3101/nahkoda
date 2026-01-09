package main

import (
	"fmt"
	"os"
	"strings"

	"nahkoda/internal/errors"
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
	executeCommand(input)
}

func executeCommand(input string) {
	// 1️⃣ PARSER
	ast, err := parser.Parse(input)
	if err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}

	// 2️⃣ SEMANTIC (STRICT)
	intent, err := semantic.Resolve(ast)
	if err != nil {
		// handle typo suggestion
		if nErr, ok := err.(*errors.NahkodaError); ok && nErr.Suggestion != "" {
			fmt.Printf("❓ %s\n", nErr.Message)
			fmt.Printf("👉 Mungkin maksud Kapten: %s? (y/n): ", nErr.Suggestion)

			var confirm string
			fmt.Scanln(&confirm)

			if strings.ToLower(confirm) == "y" {
				// Coba ganti kata pertama yang unknown dengan suggestion
				// Untuk simplicity, kita ganti kata yang mengandung typo di input string
				newInput := strings.Replace(input, ast.Unknown[0], nErr.Suggestion, 1)
				fmt.Printf("⚓ Berlayar dengan: %s\n\n", newInput)
				executeCommand(newInput)
				return
			}
		}

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
