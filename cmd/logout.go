// cmd/logout.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from SIA servers",
	Long:  `Log out from SIA servers, clearing any auth tokens.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if access token exists
		checkAccessToken()

		// delete access token
		deleteAccessToken()

		fmt.Println("successfully logged out")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
