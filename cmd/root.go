package cmd

import (
	"os"
	"strings"

	"nahkoda/internal/exec"
	"nahkoda/internal/parser"
	"nahkoda/internal/planner"
	"nahkoda/internal/semantic"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nahkoda [kalimat]",
	Short: "Nahkoda — Bahasa manusia untuk Kubernetes",
	Long: `⚓ Nahkoda — Bahasa manusia untuk Kubernetes

Contoh perintah:
  nahkoda liat kru
  nahkoda liat kru rusak
  nahkoda liat kru di geladak auth
  nahkoda liat kru bocor di geladak payment

Nahkoda menerjemahkan bahasa manusia
menjadi intent operasional Kubernetes.
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.Join(args, " ")

		// 1️⃣ PARSER
		ast, err := parser.Parse(input)
		if err != nil {
			return err
		}

		// 2️⃣ SEMANTIC (STRICT)
		intent, err := semantic.Resolve(ast)
		if err != nil {
			return err
		}

		// 3️⃣ PLANNER
		plan := planner.Build(intent)

		// 4️⃣ EXECUTOR
		if err := exec.Execute(plan); err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
