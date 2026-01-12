package main

import (
	"bufio"
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
		runREPL()
		return
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		return
	}

	input := strings.Join(os.Args[1:], " ")
	executeCommand(input)
}

func runREPL() {
	fmt.Println("⚓ Selamat datang di Geladak Nahkoda Interaktif!")
	fmt.Println("   Ketik perintah Anda (contoh: 'liat kru') atau 'keluar' untuk mengakhiri.")
	fmt.Println("")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("⚓ > ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "keluar" || input == "exit" || input == "quit" {
			fmt.Println("⚓ Kapal sedang membuang sauh. Sampai jumpa, Kapten!")
			break
		}

		executeCommand(input)
	}
}

func executeCommand(input string) {
	ast, err := parser.Parse(input)
	if err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}

	intent, err := semantic.Resolve(ast)
	if err != nil {
		if nErr, ok := err.(*errors.NahkodaError); ok && nErr.Suggestion != "" {
			fmt.Printf("❓ %s\n", nErr.Message)
			fmt.Printf("👉 Mungkin maksud Kapten: %s? (y/n): ", nErr.Suggestion)

			var confirm string
			_, _ = fmt.Scanln(&confirm)

			if strings.ToLower(confirm) == "y" {
				newInput := strings.Replace(input, ast.Unknown[0], nErr.Suggestion, 1)
				fmt.Printf("⚓ Berlayar dengan: %s\n\n", newInput)
				executeCommand(newInput)
				return
			}
		}

		fmt.Println("❌", err.Error())
		os.Exit(1)
	}

	plan := planner.Build(intent)

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
