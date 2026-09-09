package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"nahkoda/internal/completer"
	"nahkoda/internal/config"
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
	// Parse flags
	dryRun := flag.Bool("dry-run", false, "Tampilkan perintah tanpa eksekusi")
	verbose := flag.Bool("verbose", false, "Tampilkan detail eksekusi")
	flag.BoolVar(verbose, "v", false, "Tampilkan detail eksekusi")
	help := flag.Bool("help", false, "Tampilkan bantuan")
	flag.BoolVar(help, "h", false, "Tampilkan bantuan")

	flag.Usage = printHelp
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Load Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("⚠️  Gagal memuat konfigurasi: %v\n", err)
		// Invalid config must not leave runtime fields at unsafe zero values.
		cfg = config.Default()
	}

	// Dependency Injection: Initialize standard client and executor
	client := &exec.StandardKubectlClient{
		KubectlPath: cfg.KubectlPath,
		Timeout:     cfg.Timeout,
	}
	executor := exec.NewExecutor(client)
	executor.DryRun = *dryRun
	executor.Verbose = *verbose
	completer.Configure(cfg.KubectlPath, cfg.CacheTTL)

	args := flag.Args()

	if len(args) < 1 {
		runREPL(executor, cfg.DefaultNamespace, cfg.EnableSuggestions)
		return
	}

	input := strings.Join(args, " ")
	if err := processCommand(input, executor, cfg.DefaultNamespace, false); err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}
}

func runREPL(executor *exec.Executor, defaultNamespace string, enableSuggestions bool) {
	fmt.Println("⚓ Selamat datang di Anjungan Pintar Nahkoda!")
	if enableSuggestions {
		fmt.Println("   Ketik perintah Anda. Gunakan TAB untuk saran sakti.")
	} else {
		fmt.Println("   Ketik perintah Anda. Autocomplete dinonaktifkan lewat config.")
	}
	fmt.Println("   Ketik 'keluar' untuk mengakhiri pelayaran.")
	fmt.Println("")

	// Closure to pass runtime configuration to command execution.
	executeCmd := func(input string) {
		executeCommand(input, executor, defaultNamespace)
	}

	completerFunc := completer.Completer
	if !enableSuggestions {
		completerFunc = func(prompt.Document) []prompt.Suggest { return nil }
	}

	p := prompt.New(
		executeCmd,
		completerFunc,
		prompt.OptionPrefix("⚓ > "),
		prompt.OptionTitle("Nahkoda Anjungan Pintar"),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionSelectedDescriptionBGColor(prompt.Cyan),
		prompt.OptionCompletionOnDown(), // Only show suggestions when pressing Down arrow, not automatically
	)
	p.Run()
}

func executeCommand(input string, executor *exec.Executor, defaultNamespace string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	if input == "keluar" || input == "exit" || input == "quit" {
		fmt.Println("⚓ Kapal sedang membuang sauh. Sampai jumpa, Kapten!")
		os.Exit(0)
	}

	if err := processCommand(input, executor, defaultNamespace, true); err != nil {
		fmt.Println("❌", err.Error())
		// Stay in REPL, do not os.Exit(1)
	}
}

func processCommand(input string, executor *exec.Executor, defaultNamespace string, interactive bool) error {
	ast, err := parser.Parse(input)
	if err != nil {
		return err
	}

	intent, err := semantic.ResolveWithDefaultNamespace(ast, defaultNamespace)
	if err != nil {
		if nErr, ok := err.(*errors.NahkodaError); ok && nErr.Suggestion != "" {
			if !interactive {
				return fmt.Errorf("%s; mungkin maksud Kapten: %s", nErr.Message, nErr.Suggestion)
			}

			fmt.Printf("❓ %s\n", nErr.Message)
			fmt.Printf("👉 Mungkin maksud Kapten: %s? (y/n): ", nErr.Suggestion)

			var confirm string
			_, _ = fmt.Scanln(&confirm)

			if strings.ToLower(confirm) == "y" {
				newInput := strings.Replace(input, ast.Unknown[0], nErr.Suggestion, 1)
				fmt.Printf("⚓ Berlayar dengan: %s\n\n", newInput)
				return processCommand(newInput, executor, defaultNamespace, true)
			}
		}
		return err
	}

	plan := planner.Build(intent)
	return executor.Execute(plan)
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
  -h, --help       Tampilkan bantuan ini
  --dry-run        Tampilkan perintah tanpa eksekusi
  -v, --verbose    Tampilkan detail eksekusi`)
}
