package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type SetPwInput struct {
	Username       string
	Password       string
	RepeatPassword string
}

var pwInput SetPwInput

// setpwdCmd represents the set password command
var setpwdCmd = &cobra.Command{
	Use:     "setpwd",
	Short:   "Set the admin password",
	Aliases: []string{"spw"},
	Long: `
1. To set the admin password.
2. This command can be used only from the server console not a remote terminal console.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// confirm access is from server console
		confirmIfLocalHost()

		// Prompt for password
		pwInput.Password = readHiddenTextInput("Enter a strong password (min 6 chars): ")
		// prompt for repeat password
		pwInput.RepeatPassword = readVisibleTextInput("Repeat above password: ")

		// Check if passwords match
		if pwInput.Password != pwInput.RepeatPassword {
			fmt.Println("Passwords do not match.")
			fmt.Println()
			return
		}

		// set the login URL
		setpwURL := "/api/auth/set-admin-password"

		// Create the payload as JSON
		payload := map[string]string{
			"password": pwInput.Password,
		}

		// Get the request body
		reqBody := generateJSONBody(payload)

		// Create the POST request
		req := createHttpClient("POST", setpwURL, reqBody, "application/json")

		// Execute HTTP client
		res, resBody := executeHttpRequest(req)

		//Check status code
		checkResponseStatusCode(res, resBody)

		// Print the successful response
		fmt.Println("Admin password successfully set. Login to proceed.")
		fmt.Println()

	},
}

func init() {
	rootCmd.AddCommand(setpwdCmd)
}
