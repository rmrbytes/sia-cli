/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:   "sia",
	Short: "sia is a command line tool for managing your SIA servers",
	Long: `
1. sia is a command line tool for managing your SIA servers. 
2. It has commands to set up your intelligent agents.
3. Ensure that the environment variable SIA_SERVER_URL and SIA_API_KEY is set before using this tool. 
4. To set the server envar:
   - On Linux or macOS:
     export SIA_SERVER_URL=http://localhost:8080
   - On Windows:
     set SIA_SERVER_URL=https://sia.example.com
5. To set the API key:
   - On Linux or macOS:
     export SIA_API_KEY=the-access-key
   - On Windows:
     set SIA_API_KEYL=the-access-key
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the version flag is set
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Println(version)
			os.Exit(0)
		}
		cmd.Help()
	},
}

func Execute() {
	// rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add the -v or --version flag to rootCmd only (not persistent across subcommands)
	rootCmd.Flags().BoolP("version", "v", false, "Display the version of sia-cli")
}
