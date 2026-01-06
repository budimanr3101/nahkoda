package cmd

import (
	"fmt"
	"os"
	"strings"

	"nahkoda/internal/exec"
	"nahkoda/internal/parser"
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
  nahkoda liat kru bocor di geladak auth

Nahkoda akan menerjemahkan kalimat manusia
menjadi intent operasional Kubernetes.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")

		ast, err := parser.Parse(input)
		if err != nil {
			fmt.Println("❌", err.Error())
			return
		}

		intent := semantic.Resolve(ast)
		exec.Execute(intent)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
