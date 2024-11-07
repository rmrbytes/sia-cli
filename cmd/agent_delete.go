package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var agentDeleteName string

var agentDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del"},
	Short:   "Delete an existing agent",
	Long: `
Delete an existing agent`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if access token exists
		checkAccessToken()

		// Construct the delete URL
		deleteURL := fmt.Sprintf("/api/agents/%s", agentDeleteName)

		// Create Http Request
		req := createAuthHttpClient("DELETE", deleteURL, nil, "")

		// Execute HTTP client
		response, responseBody := executeHttpRequest(req)
		//Check status code
		checkResponseStatusCode(response, responseBody)
		// display success
		fmt.Println("agent successfully deleted")
	},
}

func init() {
	agentDeleteCmd.Flags().StringVarP(&agentDeleteName, "name", "n", "", "Name of the agent to delete")
	agentDeleteCmd.MarkFlagRequired("name")

	agentCmd.AddCommand(agentDeleteCmd)
}
