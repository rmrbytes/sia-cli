package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ChangePwInput struct {
	Username        string
	CurrentPassword string
	Password        string
	RepeatPassword  string
}

var changePwInput ChangePwInput

// setpwdCmd represents the set password command
var changepwdCmd = &cobra.Command{
	Use:     "changepwd",
	Short:   "Change the admin password",
	Aliases: []string{"cpw"},
	Long: `
1. To change the admin password.
2. This command can be used only from the server console not a remote terminal console.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkEnvVars()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// confirm access is from server console
		confirmIfLocalHost()

		// Prompt forcurrent password
		changePwInput.CurrentPassword = readHiddenTextInput("Enter the current admin password: ")

		// Prompt for password
		changePwInput.Password = readHiddenTextInput("Enter a new strong password (min 6 chars): ")
		// prompt for repeat password
		changePwInput.RepeatPassword = readVisibleTextInput("Repeat above password: ")

		// Check if passwords match
		if changePwInput.Password != changePwInput.RepeatPassword {
			fmt.Println("Passwords do not match.")
			fmt.Println()
			return
		}

		// set the login URL
		setpwURL := "/api/auth/update-admin-password"

		// Create the payload as JSON
		payload := map[string]string{
			"current_password": changePwInput.CurrentPassword,
			"new_password":     changePwInput.Password,
		}

		// Get the request body
		reqBody := generateJSONBody(payload)

		// Create the POST request
		req := createAuthHttpClient("POST", setpwURL, reqBody, "application/json")

		// Execute HTTP client
		res, resBody := executeHttpRequest(req)

		//Check status code
		checkResponseStatusCode(res, resBody)

		// Print the successful response
		fmt.Println("Admin password successfully changed.")
		fmt.Println()

	},
}

func init() {
	rootCmd.AddCommand(changepwdCmd)
}
