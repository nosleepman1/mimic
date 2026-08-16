package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/teranga-host/mimic/internal/parser"
)

var validateFile string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check the syntax and validity of an API contract file",
	Run: func(cmd *cobra.Command, args []string) {
		runValidate()
	},
}

func init() {
	validateCmd.Flags().StringVarP(&validateFile, "file", "f", "mimic.yaml", "Path to contract file")
}

func runValidate() {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

	if _, err := os.Stat(validateFile); os.IsNotExist(err) {
		fmt.Printf("%s File '%s' not found.\n", red("[ERROR]"), validateFile)
		os.Exit(1)
	}

	endpoints, err := parser.DetectAndParse(validateFile)
	if err != nil {
		fmt.Printf("%s Contract validation failed: %v\n", red("[INVALID]"), err)
		os.Exit(1)
	}

	fmt.Printf("%s Contract '%s' is valid (%d endpoints detected).\n\n", green("[VALID]"), cyan(validateFile), len(endpoints))
	for _, ep := range endpoints {
		fmt.Printf("  %-7s %s\n", ep.Method, ep.Path)
	}
	fmt.Println()
}
