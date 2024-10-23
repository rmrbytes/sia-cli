package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var agentViewName string

var agentViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"vi"},
	Short:   "View information about an agent",
	Long: `
View information about an agent`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if access token exists
		checkAccessToken()

		// Construct the view URL
		viewURL := fmt.Sprintf("/api/agents/%s", agentViewName)

		// Create Http Request
		req := createAuthHttpClient("GET", viewURL, nil, "")

		// Execute HTTP client
		response, responseBody := executeHttpRequest(req)
		//Check status code
		checkResponseStatusCode(response, responseBody)
		// Unmarshal response to AgentResponse
		agentResponse := unmarshalAgentResponse(responseBody)
		// Convert AgentResponse to AgentDisplay
		agentDisplay := convertAgentResponseToDisplay(agentResponse)

		// Display the agent details
		displayAgentDetails(agentDisplay)

	},
}

func init() {
	agentViewCmd.Flags().StringVarP(&agentViewName, "name", "n", "", "Name of the agent to view")
	agentViewCmd.MarkFlagRequired("name")

	agentCmd.AddCommand(agentViewCmd)
}
