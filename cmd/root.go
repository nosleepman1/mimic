package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mimic",
	Short: "Mimic - High-Performance Intelligent Mock Server and Contract Engine",
	Long: `Mimic is an ultra-fast, zero-dependency API mock server in Go.
It instantly turns OpenAPI, Swagger, Postman, Laravel routes, and YAML files
into realistic, stateful in-memory CRUD endpoints with network latency and chaos simulation.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display Mimic version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Mimic API Mock Engine v1.0.0 (Go runtime)")
	},
}
