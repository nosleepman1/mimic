package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a starter mimic.yaml template file",
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func runInit() {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

	targetFile := "mimic.yaml"
	if _, err := os.Stat(targetFile); err == nil {
		fmt.Printf("[WARN] '%s' already exists in current directory.\n", targetFile)
		return
	}

	sampleContent := `# mimic.yaml - Declarative API Mock Definition
# Run: mimic serve -f mimic.yaml --port=8080 --latency=200ms

endpoints:
  # 1. Stateful CRUD Users Resource
  - path: /api/v1/users
    method: GET
    summary: "List all users"
    resource: users
    is_crud: true

  - path: /api/v1/users
    method: POST
    summary: "Create a new user"
    resource: users
    is_crud: true
    status: 201

  - path: /api/v1/users/{id}
    method: GET
    summary: "Retrieve user by ID"
    resource: users
    is_crud: true

  - path: /api/v1/users/{id}
    method: PUT
    summary: "Update user by ID"
    resource: users
    is_crud: true

  - path: /api/v1/users/{id}
    method: DELETE
    summary: "Delete user"
    resource: users
    is_crud: true
    status: 204

  # 2. Custom Payload / Static Response
  - path: /api/v1/auth/login
    method: POST
    summary: "Authenticate user"
    status: 200
    response:
      token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
      token_type: "Bearer"
      expires_in: 86400

  # 3. Payments with Localized XOF Currency & Senegalese Phone Numbers
  - path: /api/v1/payments/mobile-money
    method: POST
    summary: "Initiate Wave or Orange Money transaction"
    status: 200
    schema:
      type: object
      properties:
        transaction_id:
          type: uuid
        customer_name:
          type: name
        phone_number:
          type: phone_sn
        amount:
          type: currency_xof
        status:
          type: string
          enum: ["SUCCESS", "PENDING"]
        created_at:
          type: date_iso
`

	if err := os.WriteFile(targetFile, []byte(sampleContent), 0644); err != nil {
		fmt.Printf("[ERROR] Failed to write %s: %v\n", targetFile, err)
		return
	}

	fmt.Printf("%s Generated '%s' successfully.\n", green("[SUCCESS]"), targetFile)
	fmt.Printf("Start your mock server with:\n")
	fmt.Printf("  %s\n\n", cyan("mimic serve -f mimic.yaml --port=8080 --latency=200ms"))
}
