package cmd

import (
	"github.com/spf13/cobra"
)

// agentCmd represents the agent parent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agents with many subcommands",
	Long: `
Manage agents with subcommands like ls, create, view, pull, push, and delete.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
