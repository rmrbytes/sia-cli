package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Array to store chat history
var chatAgentName string

// Array to store chat history
var messages []ChatMessage

// chatCmd represents the chat command
var agtChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session",
	Long:  `The chat session allows you to interact with the LLM. Type 'q' to quit.`,
	Run:   startChatLoop,
}

func init() {
	agentCmd.AddCommand(agtChatCmd) // Add the chat command to the root command

	agtChatCmd.Flags().StringVarP(&chatAgentName, "name", "n", "", "Specify the agent name (required)")
	agtChatCmd.MarkFlagRequired("name")
}

// startChatLoop starts an interactive chat loop
func startChatLoop(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println("Starting chat session. Type 'q' to quit.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Display a prompt
		fmt.Print("You   : ")

		// Read user input
		scanned := scanner.Scan()
		if !scanned {
			fmt.Println("\n[Error]: Unable to read input. Exiting.")
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Handle quit command
		if input == "q" {
			fmt.Println()
			fmt.Println("Exiting chat session.")
			fmt.Println()
			break
		}

		messages = append(messages, ChatMessage{Role: "user", Content: input})
		fmt.Println("Agent : ... ")
		// Call a function to handle the chat input and get a response
		response := sendChatPrompt(chatAgentName, input, messages)
		fmt.Print("\033[F\033[K")
		fmt.Println("Agent :", response.Content)
	}

}

func sendChatPrompt(agentName, prompt string, messages []ChatMessage) ChatResponse {

	// set the chat URL & method
	chatURL := fmt.Sprintf("/api/chat/%s", agentName)
	method := "POST"
	// Prepare the request payload
	payload := ChatRequest{
		Prompt:   prompt,
		Messages: messages,
	}

	// Get the request body
	reqBody := generateJSONBody(payload)

	// Create the POST request
	req := createHttpClient(method, chatURL, reqBody, "application/json")

	// Execute HTTP client
	res, resBody := executeHttpRequest(req)

	//Check status code
	checkResponseStatusCode(res, resBody)

	// Unmarshal response to ChatResponse
	chatResponse := unmarshalChatResponse(resBody)

	return chatResponse

}
