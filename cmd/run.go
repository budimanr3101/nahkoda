package cmd

import (
	"strings"

	"nahkoda/internal/exec"
	"nahkoda/internal/parser"
	"nahkoda/internal/semantic"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [kalimat]",
	Short: "Jalankan perintah Nahkoda",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")

		ast := parser.Parse(input)
		intent := semantic.Resolve(ast)
		exec.Execute(intent)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
