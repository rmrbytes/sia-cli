package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var agentPullName string

var agentPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Download info of agent in YAML format",
	Long: `
Download info of agent in YAML format so that it may be edited and pushed to update the server.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if access token exists
		checkAccessToken()
		// Construct the pull URL
		pullURL := fmt.Sprintf("/api/agents/%s", agentPullName)

		// Create Http Request
		req := createAuthHttpClient("GET", pullURL, nil, "")

		// Execute HTTP client
		response, responseBody := executeHttpRequest(req)
		//Check status code
		checkResponseStatusCode(response, responseBody)
		// Unmarshal response to AgentResponse
		agentResponse := unmarshalAgentResponse(responseBody)
		// Unmarshal Response to AgentInputYaml
		agentInput := unmarshalAgentInputYaml(responseBody)

		// Add DeletedFiles from Existing Files
		addDeletedFiles(&agentInput, agentResponse)

		// Add Sample NewFiles Data
		addSampleNewFiles(&agentInput)

		// Marshal to YAML with Comments
		yamlWithComments := addCommentsToYaml(agentInput)

		// Step 6: Save YAML to File
		filename := fmt.Sprintf("%s.yaml", agentPullName)
		saveYamlToFile(yamlWithComments, filename)

		fmt.Printf("Agent data has been download as %s in cwd", filename)
	},
}

func init() {
	agentPullCmd.Flags().StringVarP(&agentPullName, "name", "n", "", "Name of the agent to pull")
	agentPullCmd.MarkFlagRequired("name")

	agentCmd.AddCommand(agentPullCmd)
}
