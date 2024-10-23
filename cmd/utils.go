// cmd/utils.go

package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func readAgentYamlFile(filePath string) AgentInputYaml {
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		handleErr(err, "Failed to read YAML file")
	}

	var agentInput AgentInputYaml
	err = yaml.Unmarshal(yamlData, &agentInput)
	if err != nil {
		handleErr(err, "Failed to decode YAML data")
	}

	return agentInput
}

func convertAgentInputToPushRequest(input AgentInputYaml) AgentPushRequest {
	var files []FileDetail
	for _, newFile := range input.NewFiles {
		filename := filepath.Base(newFile.Filepath)
		files = append(files, FileDetail{
			Filename: filename,
			Meta:     newFile.Meta,
		})
	}

	return AgentPushRequest{
		Name:             input.Name,
		Instructions:     input.Instructions,
		WelcomeMessage:   input.WelcomeMessage,
		SuggestedPrompts: input.SuggestedPrompts,
		DeletedFiles:     input.DeletedFiles,
		Files:            files,
	}
}

func createMultipartForm(agentRequest AgentPushRequest, input AgentInputYaml) (io.Reader, string) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the JSON fields
	err := writer.WriteField("name", agentRequest.Name)
	if err != nil {
		handleErr(err, "Failed to add name to form")
	}
	err = writer.WriteField("instructions", agentRequest.Instructions)
	if err != nil {
		handleErr(err, "Failed to add instructions to form")
	}
	err = writer.WriteField("welcome_message", agentRequest.WelcomeMessage)
	if err != nil {
		handleErr(err, "Failed to add welcome_message to form")
	}
	for _, prompt := range agentRequest.SuggestedPrompts {
		err = writer.WriteField("suggested_prompts", prompt)
		if err != nil {
			handleErr(err, "Failed to add suggested prompt to form")
		}
	}
	for _, deletedFile := range agentRequest.DeletedFiles {
		err = writer.WriteField("deleted_files", deletedFile)
		if err != nil {
			handleErr(err, "Failed to add deleted file to form")
		}
	}

	// Add metadata for each new file under "files" field
	var filesArray []FileDetail

	for _, newFile := range input.NewFiles {
		// Create a FileDetail instance using newFile's data
		fileDetail := FileDetail{
			Filename: filepath.Base(newFile.Filepath),
			Meta:     newFile.Meta,
		}
		filesArray = append(filesArray, fileDetail)
	}

	// Marshal the array of FileDetail structs to JSON
	filesMetadataJSON, err := json.Marshal(filesArray)
	if err != nil {
		handleErr(err, "Failed to marshal FileDetail array to JSON")
	}

	// Add the JSON-encoded files metadata to the form under the "files" field
	err = writer.WriteField("files", string(filesMetadataJSON))
	if err != nil {
		handleErr(err, "Failed to add files metadata to form")
	}

	// Add the actual files to be uploaded under "new_files" field
	for _, newFile := range input.NewFiles {
		resolvedPath := resolvePath(newFile.Filepath)
		file, err := os.Open(resolvedPath)
		if err != nil {
			handleErr(err, fmt.Sprintf("Failed to open file %s", newFile.Filepath))
		}
		defer file.Close()

		// Create a form file for the multipart data
		part, err := writer.CreateFormFile("new_files", filepath.Base(newFile.Filepath))
		if err != nil {
			handleErr(err, fmt.Sprintf("Failed to create form file for %s", newFile.Filepath))
		}

		_, err = io.Copy(part, file)
		if err != nil {
			handleErr(err, fmt.Sprintf("Failed to copy file data for %s", newFile.Filepath))
		}
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		handleErr(err, "Failed to close the multipart writer")
	}

	return &requestBody, writer.FormDataContentType()
}

func createHttpClient(method, url string, body io.Reader, contentType string) *http.Request {
	serverUrl := os.Getenv(("SIA_SERVER_URL"))
	fullUrl := fmt.Sprintf("%s%s", serverUrl, url)
	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		handleErr(err, "Failed to create HTTP request")
	}

	req.Header.Set("Content-Type", contentType)
	// get API key
	siaAPIKey := os.Getenv("SIA_API_KEY")
	req.Header.Set("X-Requested-With", siaAPIKey)
	return req
}

func createAuthHttpClient(method string, url string, body io.Reader, contentType string) *http.Request {
	req := createHttpClient(method, url, body, contentType)
	// add cookie to header
	accessToken := checkAccessToken()
	if accessToken == nil {
		err := errors.New("access token is not found")
		handleErr(err, "No access token found")
	}
	req.AddCookie(&http.Cookie{Name: "access_token", Value: strings.TrimSpace(string(accessToken))})
	return req
}

func executeHttpRequest(req *http.Request) (*http.Response, []byte) {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		handleErr(err, "Failed to execute HTTP request")
	}

	responseBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		handleErr(err, "Failed to read response body")
	}

	return resp, responseBody
}

func unmarshalAgentResponse(responseBody []byte) AgentResponse {
	var agentResponse AgentResponse
	err := json.Unmarshal(responseBody, &agentResponse)
	if err != nil {
		handleErr(err, "Failed to unmarshal response body")
	}

	return agentResponse
}

func convertAgentResponseToDisplay(response AgentResponse) AgentDisplay {
	return AgentDisplay{
		Name:             response.Name,
		WelcomeMessage:   response.WelcomeMessage,
		Instructions:     response.Instructions,
		SuggestedPrompts: response.SuggestedPrompts,
		Files:            response.Files,
		CreatedOn:        formatTimestamp(response.CreatedOn),
		UpdatedOn:        formatTimestamp(response.UpdatedOn),
	}
}

func displayAgentDetails(agentDisplay AgentDisplay) {
	yamlData, err := yaml.Marshal(agentDisplay)
	if err != nil {
		handleErr(err, "Failed to marshal agent details for display")
	}

	fmt.Println(string(yamlData))
}

func unmarshalAgentsListResponse(responseBody []byte) []AgentResponse {
	var agentsList []AgentResponse
	err := json.Unmarshal(responseBody, &agentsList)
	if err != nil {
		fmt.Println(err)
		handleErr(err, "Failed to unmarshal response body")
	}

	return agentsList
}

func convertAgentsListToDisplay(agentsList []AgentResponse) []AgentSummaryDisplay {
	var displayList []AgentSummaryDisplay

	for i, agent := range agentsList {
		displayList = append(displayList, AgentSummaryDisplay{
			Srno:             i + 1,
			Name:             agent.Name,
			FileCount:        len(agent.Files),
			EmbeddingsStatus: agent.EmbeddingsStatus,
			CreatedOn:        formatTimestamp(agent.CreatedOn),
			UpdatedOn:        formatTimestamp(agent.UpdatedOn),
		})
	}

	return displayList
}

func displayAgentsTable(agents []AgentSummaryDisplay) {
	// Print the header row
	headerFormat := "%-5s %-20s %-8s %-9s %-10s %-10s\n"
	fmt.Printf(headerFormat, "SRNO", "NAME", "# FILES", "E STATUS", "CREATED ON", "UPDATED ON")

	// Print a separator row for better readability
	line := strings.Repeat("-", 67)
	fmt.Println(line)

	// Print each agent's details
	rowFormat := "%-5d %-20s %-8d %-9s %-10s %-10s\n"
	for _, agent := range agents {
		fmt.Printf(rowFormat, agent.Srno, agent.Name, agent.FileCount, agent.EmbeddingsStatus, agent.CreatedOn, agent.UpdatedOn)
	}
	fmt.Println()
}

func unmarshalAgentInputYaml(responseBody []byte) AgentInputYaml {
	var agentResponse AgentResponse
	err := json.Unmarshal(responseBody, &agentResponse)
	if err != nil {
		handleErr(err, "Failed to unmarshal response body")
	}

	// Convert to AgentInputYaml (fill out the necessary fields)
	return AgentInputYaml{
		Name:             agentResponse.Name,
		Instructions:     agentResponse.Instructions,
		WelcomeMessage:   agentResponse.WelcomeMessage,
		SuggestedPrompts: agentResponse.SuggestedPrompts,
		DeletedFiles:     []string{},        // Will be filled later
		NewFiles:         []NewFileDetail{}, // Will be filled later
	}
}

func addDeletedFiles(agentInput *AgentInputYaml, agentResponse AgentResponse) {
	for _, file := range agentResponse.Files {
		agentInput.DeletedFiles = append(agentInput.DeletedFiles, file.Filename)
	}
}

func addSampleNewFiles(agentInput *AgentInputYaml) {
	agentInput.NewFiles = []NewFileDetail{
		{
			Filepath: "~/docs/document1.pdf",
			Meta: Meta{
				SplitBy:        "word",
				SplitLength:    200,
				SplitOverlap:   20,
				SplitThreshold: 0,
			},
		},
		{
			Filepath: "../files/file1.txt",
			Meta: Meta{
				SplitBy:        "paragraph",
				SplitLength:    100,
				SplitOverlap:   10,
				SplitThreshold: 0,
			},
		},
	}
}

func addCommentsToYaml(agentInput AgentInputYaml) string {
	// Marshal the struct to YAML
	yamlData, err := yaml.Marshal(agentInput)
	if err != nil {
		handleErr(err, "Failed to marshal AgentInputYaml")
	}

	// Add comments to specific fields in the YAML
	lines := strings.Split(string(yamlData), "\n")
	var builder strings.Builder

	inDeletedFilesSection := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(line, "name:") {
			builder.WriteString(fmt.Sprintf("%s # this cannot be changed\n", line))
			inDeletedFilesSection = false
		} else if strings.HasPrefix(line, "suggested_prompts:") {
			builder.WriteString(fmt.Sprintf("%s # a max of 3 prompts can be given\n", line))
			inDeletedFilesSection = false
		} else if strings.HasPrefix(line, "deleted_files:") {
			builder.WriteString(fmt.Sprintf("%s # uncomment the files you wish to delete\n", line))
			inDeletedFilesSection = true
		} else if inDeletedFilesSection && strings.HasPrefix(trimmedLine, "- ") {
			index := strings.Index(line, "-")
			commentedLine := fmt.Sprintf("%s# %s", line[:index], line[index:])
			builder.WriteString(commentedLine + "\n")
		} else if strings.HasPrefix(line, "new_files:") {
			builder.WriteString(fmt.Sprintf("%s # change below template as required\n", line))
			inDeletedFilesSection = false
		} else if strings.HasPrefix(line, "- ") && strings.HasPrefix(strings.TrimSpace(line), "#") {
			// Leave existing comments as-is
			builder.WriteString(fmt.Sprintf("%s\n", line))
		} else {
			builder.WriteString(line + "\n")
			inDeletedFilesSection = false
		}
	}

	return builder.String()
}

func saveYamlToFile(yamlData string, filename string) {
	err := os.WriteFile(filename, []byte(yamlData), 0644)
	if err != nil {
		handleErr(err, "Failed to write YAML to file")
	}
}

// Converts an int64 timestamp to a formatted date string.
func formatTimestamp(timestamp int64) string {
	// Convert Unix timestamp to time.Time
	t := time.Unix(timestamp, 0)

	// Format the time object to "DD-MMM-YY"
	return t.Format("15-Jan-06")
}

func generateJSONBody(payload interface{}) *bytes.Buffer {
	// Marshal the input payload to JSON
	requestBody, err := json.Marshal(payload)
	if err != nil {
		handleErr(err, "Failed to marshal request payload")
	}
	// Wrap the marshaled JSON in a bytes.Buffer and return it
	return bytes.NewBuffer(requestBody)
}

func checkResponseStatusCode(res *http.Response, body []byte) {
	// Check the response status code for success (200-299)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// Parse the response body to extract the "detail" field
		var errorResponse struct {
			Detail string `json:"detail"`
		}
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			handleErr(err, "Server error")
		}
		handleErr(err, errorResponse.Detail)
		return
	}
}

func retrieveTokenAndSave(res *http.Response) {
	var accessToken string
	// retrieve access token from cookies
	for _, cookie := range res.Cookies() {
		if cookie.Name == "access_token" {
			accessToken = cookie.Value
			break
		}
	}
	// check if blank
	if accessToken == "" {
		err := errors.New("access token is blank")
		handleErr(err, "No access token found")
	}
	// save it in local directory
	saveAccessToken(accessToken)
}
