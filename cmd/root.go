package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nahkoda",
	Short: "Nahkoda - berbicara ke Kubernetes dengan bahasa manusia",
	Long:  "Nahkoda adalah CLI untuk mengendalikan Kubernetes menggunakan bahasa Indonesia.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
