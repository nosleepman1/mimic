package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/teranga-host/mimic/internal/parser"
	"github.com/teranga-host/mimic/internal/server"
)

var (
	serveFile    string
	servePort    int
	serveHost    string
	serveLatency string
	serveJitter  string
	serveChaos   float64
	serveCORS    bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the mock HTTP server from an API contract file",
	Long: `Parses an OpenAPI, Swagger, Postman, Laravel routes/api.php, or mimic.yaml file
and launches a high-performance in-memory mock server with realistic dynamic data generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServe()
	},
}

func init() {
	serveCmd.Flags().StringVarP(&serveFile, "file", "f", "mimic.yaml", "Path to API contract file (OpenAPI, Postman, Laravel api.php, or YAML)")
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Port to listen on")
	serveCmd.Flags().StringVar(&serveHost, "host", "0.0.0.0", "Host address to bind to")
	serveCmd.Flags().StringVar(&serveLatency, "latency", "", "Simulated network latency (e.g. '200ms', '1s')")
	serveCmd.Flags().StringVar(&serveJitter, "jitter", "", "Simulated network jitter (e.g. '50ms')")
	serveCmd.Flags().Float64Var(&serveChaos, "chaos", 0.0, "Chaos 500 error rate (0.0 to 1.0, e.g. 0.05 for 5%)")
	serveCmd.Flags().BoolVar(&serveCORS, "cors", true, "Enable CORS headers for frontend/mobile clients")
}

func runServe() {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

	if _, err := os.Stat(serveFile); os.IsNotExist(err) {
		fmt.Printf("%s Contract file '%s' not found.\n", red("[ERROR]"), serveFile)
		fmt.Printf("Run '%s' to create a sample contract file.\n\n", cyan("mimic init"))
		os.Exit(1)
	}

	fmt.Printf("[INFO] Ingesting contract from: %s\n", cyan(serveFile))
	endpoints, err := parser.DetectAndParse(serveFile)
	if err != nil {
		fmt.Printf("%s Failed to parse contract: %v\n", red("[ERROR]"), err)
		os.Exit(1)
	}

	if len(endpoints) == 0 {
		fmt.Printf("%s No endpoints detected in file '%s'.\n", red("[WARN]"), serveFile)
		os.Exit(1)
	}

	var latencyDuration, jitterDuration time.Duration
	if serveLatency != "" {
		latencyDuration, _ = time.ParseDuration(serveLatency)
	}
	if serveJitter != "" {
		jitterDuration, _ = time.ParseDuration(serveJitter)
	}

	cfg := server.Config{
		Port:       servePort,
		Host:       serveHost,
		Latency:    latencyDuration,
		Jitter:     jitterDuration,
		ChaosRate:  serveChaos,
		EnableCORS: serveCORS,
	}

	mockServer := server.NewServer(cfg, endpoints)
	if err := mockServer.Start(); err != nil {
		fmt.Printf("%s Server error: %v\n", red("[FATAL]"), err)
		os.Exit(1)
	}
}
