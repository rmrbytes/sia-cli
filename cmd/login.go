package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Password flag variable
var loginPassword string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into SIA servers",
	Long: `
1. To log into your SIA servers as an admin.
2. If the password is not part of the command, the app will prompt you for it.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Check if the password flag was provided
		if loginPassword == "" {
			// otherwise prompt the user
			loginPassword = readHiddenTextInput("Enter Admin Password:")
		}

		// set the login URL
		loginURL := "/api/auth/login"

		// Create the payload as JSON
		payload := map[string]string{
			"password": loginPassword,
		}

		// Get the request body
		reqBody := generateJSONBody(payload)

		// Create the POST request
		req := createHttpClient("POST", loginURL, reqBody, "application/json")

		// Execute HTTP client
		res, resBody := executeHttpRequest(req)

		//Check status code
		checkResponseStatusCode(res, resBody)

		// retrieve token from cookie
		retrieveTokenAndSave(res)

		// Print the successful response
		fmt.Println("You are logged in.")
		fmt.Println()
	},
}

func init() {
	// Add the password flag
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for admin login")

	// Add the `login` command to the root command
	rootCmd.AddCommand(loginCmd)
}
