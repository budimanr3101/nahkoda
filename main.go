package main

import (
	"fmt"
	"os"
	"strings"

	"nahkoda/internal/completer"
	"nahkoda/internal/errors"
	"nahkoda/internal/exec"
	"nahkoda/internal/parser"
	"nahkoda/internal/planner"
	"nahkoda/internal/semantic"

	"github.com/c-bata/go-prompt"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
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
	if err := processCommand(input); err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}
}

func runREPL() {
	fmt.Println("⚓ Selamat datang di Anjungan Pintar Nahkoda!")
	fmt.Println("   Ketik perintah Anda. Gunakan TAB untuk saran sakti.")
	fmt.Println("   Ketik 'keluar' untuk mengakhiri pelayaran.")
	fmt.Println("")

	p := prompt.New(
		executeCommand,
		completer.Completer,
		prompt.OptionPrefix("⚓ > "),
		prompt.OptionTitle("Nahkoda Anjungan Pintar"),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionSelectedDescriptionBGColor(prompt.Cyan),
	)
	p.Run()
}

func executeCommand(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	if input == "keluar" || input == "exit" || input == "quit" {
		fmt.Println("⚓ Kapal sedang membuang sauh. Sampai jumpa, Kapten!")
		os.Exit(0)
	}

	if err := processCommand(input); err != nil {
		fmt.Println("❌", err.Error())
		// Stay in REPL, do not os.Exit(1)
	}
}

func processCommand(input string) error {
	ast, err := parser.Parse(input)
	if err != nil {
		return err
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
				return processCommand(newInput)
			}
		}
		return err
	}

	plan := planner.Build(intent)
	return exec.Execute(plan)
}

func printHelp() {
	helpText := fmt.Sprintf(`⚓ Nahkoda v%s — Bahasa manusia untuk Kubernetes
(Versi: %s, Commit: %s)

Cara pakai:
  nahkoda [kalimat perintah]`, Version, Version, Commit)

	fmt.Println(helpText)
	fmt.Println(`Contoh:
  nahkoda liat kru
  nahkoda liat kru rusak
  nahkoda liat kru di geladak auth
  nahkoda baca jurnal healthy-pod-1
  nahkoda masuk healthy-pod-1
  nahkoda liat berita

Opsi:
  -h, --help    Tampilkan bantuan ini`)
}
