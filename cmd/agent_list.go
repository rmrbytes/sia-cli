package cmd

import (
	"github.com/spf13/cobra"
)

// agentListCmd represents the subcommand for listing agents
var agentListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all agents",
	Long:    "List all agents on SIA servers",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if access token exists
		checkAccessToken()

		// Construct the view URL
		listURL := "/api/agents/"

		// Create Http Request
		req := createAuthHttpClient("GET", listURL, nil, "")

		// Execute HTTP client
		response, responseBody := executeHttpRequest(req)

		//Check status code
		checkResponseStatusCode(response, responseBody)

		// Unmarshal response to AgentsListResponse
		agentsList := unmarshalAgentsListResponse(responseBody)

		// Convert AgentsListResponse to AgentSummaryDisplay list
		agentsDisplayList := convertAgentsListToDisplay(agentsList)

		// Display the list in a table format
		displayAgentsTable(agentsDisplayList)

	},
}

func init() {
	// Add the `list` subcommand to `agent` parent command
	agentCmd.AddCommand(agentListCmd)
}
