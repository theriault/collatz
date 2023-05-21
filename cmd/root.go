package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "collatz",
		Short: "A CLI tool for generating information about the Collatz Conjecture.",
	}
)

func Execute() error {
	rootCmd.AddCommand(timeCmd)
	rootCmd.AddCommand(maxCmd)
	rootCmd.AddCommand(ratiosCmd)
	rootCmd.AddCommand(compareCmd)
	return rootCmd.Execute()
}

func init() {

}
