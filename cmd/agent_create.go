package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// agentListCmd represents the subcommand for listing agents
var agentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "To download a create template",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if access token exists
		checkAccessToken()

		// Step 1: Create the multiline YAML content
		yamlContent := `
name: agent-name # A meaningful name with letters, digits, hyphen, underscore, no blanks
instructions: |
  This is a sample instruction for the agent. It can be multiline.
  
  Edit it accordingly.
welcome_message: Welcome to the Agent!
suggested_prompts: # A max of 3 prompts can be given
  - What can you do?
  - How do I use this agent?
  - Tell me something interesting.
new_files:
  - filepath: "~/docs/document1.pdf" # absolute path
    meta:
      split_by: "sentence"
      split_length: 4
      split_overlap: 1
      split_threshold: 0
  - filepath: "../files/file1.txt" # relative path to cwd
    meta:
      split_by: "word" 
      split_length: 200 # defaults will be used for missing meta
`

		// Save the YAML file in the current working directory
		saveYamlToFile(yamlContent, "create-agent.yaml")

		fmt.Printf("Template YAML file for new agent has been downloaded to cwd.\n")
		fmt.Println()
	},
}

func init() {
	// Add the `create` subcommand to `agent` parent command
	agentCmd.AddCommand(agentCreateCmd)
}
