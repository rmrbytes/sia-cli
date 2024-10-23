package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/term"
)

const (
	TokenDir      = ".sia"
	TokenFilename = ".access_token"
)

// handleErr to handle errors for non-command functions
func handleErr(err error, msg string) {
	if msg != "" {
		fmt.Printf("Error: %s\n", msg)
	} else if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Unknown error\n")
	}
	os.Exit(1)
}

// Generic must function that takes a value, an error, and a custom message
func must[T any](value T, err error, msg string) T {
	if err != nil {
		handleErr(err, msg)
	}
	return value
}

// check EnvVars
func checkEnvVars() {
	siaServerURL := os.Getenv("SIA_SERVER_URL")
	siaAPIKey := os.Getenv("SIA_API_KEY")
	if siaServerURL == "" || siaAPIKey == "" {
		handleErr(fmt.Errorf("SIA_SERVER_URL and SIA_API_KEY must be set before using this CLI"), "")
	}
}

// to check if url is localhost
func confirmIfLocalHost() {
	// Define a regex pattern to match "http://localhost" on any port
	urlPattern := `^(http|https)://localhost(:\d+)?$`

	var err error
	// get the url
	serverURL := os.Getenv("SIA_SERVER_URL")
	if serverURL == "" {
		err = errors.New("SIA_SERVER_URL has not been set")
		handleErr(err, "")
	}

	// check if matches
	matched, err := regexp.MatchString(urlPattern, serverURL)
	if err != nil {
		err = errors.New("failed to validate server URL")
	}
	// return results
	if !matched {
		err = errors.New("Access is permitted only from server console")
	}
	return
}

// save token
func saveAccessToken(accessToken string) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleErr(err, "failed to retrieve user home directory")
	}

	// Define the path to save the access token
	tokenFilePath := filepath.Join(homeDir, TokenDir, TokenFilename)

	// Ensure the `.sia` directory exists
	siaDir := filepath.Join(homeDir, ".sia")
	if _, err := os.Stat(siaDir); os.IsNotExist(err) {
		err := os.Mkdir(siaDir, 0700)
		if err != nil {
			message := fmt.Sprintf("failed to created directory %s:", siaDir)
			handleErr(err, message)
		}
	}

	// Save the access token to the file
	err = os.WriteFile(tokenFilePath, []byte(accessToken), 0600)
	if err != nil {
		handleErr(err, "failed to save access token")
	}
}

// delete token
func deleteAccessToken() {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleErr(err, "failed to retrieve user home directory")
	}

	// Define the path to save the access token
	tokenFilePath := filepath.Join(homeDir, TokenDir, TokenFilename)

	// Ensure the `.sia` directory exists
	siaDir := filepath.Join(homeDir, ".sia")
	if _, err := os.Stat(siaDir); os.IsNotExist(err) {
		err := os.Mkdir(siaDir, 0700)
		if err != nil {
			message := fmt.Sprintf("failed to created directory %s:", siaDir)
			handleErr(err, message)
		}
	}

	// delete the the file
	err = os.Remove(tokenFilePath)
	if err != nil {
		handleErr(err, "failed to delete access token")
	}
}

func checkAccessToken() []byte {
	homeDir, _ := os.UserHomeDir()
	tokenFilePath := filepath.Join(homeDir, TokenDir, TokenFilename)
	accessToken, err := os.ReadFile(tokenFilePath)
	if err != nil || accessToken == nil {
		handleErr(err, "Login required. Use 'sia login'")
	}
	return accessToken
}

// to read hidden Input
func readHiddenTextInput(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		handleErr(err, "text entry error")
	}
	fmt.Println() // Newline after input is required
	return strings.TrimSpace(string(bytePassword))
}

// readVisiblePassword reads a password with visible input
func readVisibleTextInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		handleErr(err, "text entry error")
	}
	return strings.TrimSpace(line)
}

// resolve path
func resolvePath(path string) string {
	// Check if the path starts with "~", indicating the home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			handleErr(err, "failed to get home directory")
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Convert to absolute path for relative paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		handleErr(err, "failed to get absolute path")
	}

	return absPath
}

// delete file
func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	return err
}
