package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var agentPushName string
var agentPushFilePath string
var agentPushAction string

var agentPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a new or updated agent info in YAML format to the backend",
	Long: `
Push a new or updated agent info in YAML format to the backend. 
	
1. Note that the YAML format to "create" a new agent and that of an "update" is same and the PUSH subcommand is used for both.
2. Hence the need to specify action 
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if access token exists
		checkAccessToken()

		// Validate action flag
		if agentPushAction != "create" && agentPushAction != "update" {
			fmt.Println("Error: Action must be either 'create' or 'update'.")
			return
		}

		// Validate that file exists
		if _, err := os.Stat(agentPushFilePath); os.IsNotExist(err) {
			fmt.Printf("Error: File %s not found.\n", agentPushFilePath)
			return
		}

		// Step 1: Read YAML file
		agentInput := readAgentYamlFile(agentPushFilePath)

		// Step 2: Convert AgentInputYaml to AgentPushRequest
		agentRequest := convertAgentInputToPushRequest(agentInput)

		// Step 3: Create multipart form with files and JSON data
		requestBody, contentType := createMultipartForm(agentRequest, agentInput)

		// Step 4: Create HTTP client and request
		var method string
		var url string
		if agentPushAction == "create" {
			method = "POST"
			url = "/api/agents/"
		} else {
			method = "PUT"
			url = fmt.Sprintf("/api/agents/%s", agentPushName)
		}
		req := createAuthHttpClient(method, url, requestBody, contentType)

		// Execute HTTP client
		response, responseBody := executeHttpRequest(req)

		//Check status code
		checkResponseStatusCode(response, responseBody)

		// Unmarshal response to AgentResponse
		agentResponse := unmarshalAgentResponse(responseBody)

		// Convert AgentResponse to AgentDisplay
		agentDisplay := convertAgentResponseToDisplay(agentResponse)

		// display on terminal
		fmt.Println("Agent has been updated")
		fmt.Println("----------------------")
		// Display the agent details
		displayAgentDetails(agentDisplay)
		fmt.Println()
		err := deleteFile(agentPushFilePath)
		if err != nil {
			fmt.Printf("%s could not be deleted", agentPushFilePath)
		} else {
			fmt.Printf("%s has been deleted", agentPushFilePath)
		}
		fmt.Println()

	},
}

func init() {
	agentPushCmd.Flags().StringVarP(&agentPushName, "name", "n", "", "Name of the agent")
	agentPushCmd.Flags().StringVarP(&agentPushFilePath, "file", "f", "", "Path to the YAML file")
	agentPushCmd.Flags().StringVarP(&agentPushAction, "action", "a", "", "Action to perform: create or update")

	agentPushCmd.MarkFlagRequired("name")
	agentPushCmd.MarkFlagRequired("file")
	agentPushCmd.MarkFlagRequired("action")

	agentCmd.AddCommand(agentPushCmd)
}
