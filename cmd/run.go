package cmd

import (
	"fmt"
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

		ast, err := parser.Parse(input)
		if err != nil {
			fmt.Println("❌", err.Error())
			return
		}

		intent := semantic.Resolve(ast)
		exec.Execute(intent)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
