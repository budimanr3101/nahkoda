package cmd

import (
	"fmt"
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
	Short: "Nahkoda — Natural language CLI for Kubernetes",
	Long: `⚓ Nahkoda — Bahasa manusia untuk Kubernetes

Contoh perintah:
  nahkoda liat kru
  nahkoda liat kru rusak
  nahkoda liat kru di geladak auth
  nahkoda liat kru bocor di geladak payment

Analogi:
  kru       → pod
  geladak   → namespace
  rusak     → status!=Running
  bocor     → OOMKilled
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.Join(args, " ")

		ast, err := parser.Parse(input)
		if err != nil {
			return err
		}

		intent := semantic.Resolve(ast)
		plan := planner.Build(intent)

		return exec.Execute(plan)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌", err.Error())
		os.Exit(1)
	}
}
