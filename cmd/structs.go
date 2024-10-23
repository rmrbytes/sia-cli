package cmd

type Meta struct {
	SplitBy        string `json:"split_by" yaml:"split_by"`
	SplitLength    int    `json:"split_length" yaml:"split_length"`
	SplitOverlap   int    `json:"split_overlap" yaml:"split_overlap"`
	SplitThreshold int    `json:"split_threshold" yaml:"split_threshold"`
}

type FileDetail struct {
	Filename string `json:"filename" yaml:"filename"`
	Meta     Meta   `json:"meta" yaml:"meta"`
}

type NewFileDetail struct {
	Filepath string `json:"filepath" yaml:"filepath"`
	Meta     Meta   `json:"meta" yaml:"meta"`
}

type AgentDisplay struct {
	Name             string
	WelcomeMessage   string
	Instructions     string
	SuggestedPrompts []string
	Files            []FileDetail
	CreatedOn        string
	UpdatedOn        string
}

type AgentSummaryDisplay struct {
	Srno             int
	Name             string
	FileCount        int
	EmbeddingsStatus string
	CreatedOn        string
	UpdatedOn        string
}

type AgentInputYaml struct {
	Name             string          `yaml:"name"`
	Instructions     string          `yaml:"instructions"`
	WelcomeMessage   string          `yaml:"welcome_message"`
	SuggestedPrompts []string        `yaml:"suggested_prompts"`
	DeletedFiles     []string        `yaml:"deleted_files"`
	NewFiles         []NewFileDetail `yaml:"new_files"`
}

type AgentPushRequest struct {
	Name             string          `json:"name"`
	Instructions     string          `json:"instructions"`
	WelcomeMessage   string          `json:"welcome_message"`
	SuggestedPrompts []string        `json:"suggested_prompts"`
	Files            []FileDetail    `json:"files"`
	DeletedFiles     []string        `json:"deleted_files"`
	NewFiles         []NewFileDetail `json:"new_files"`
}

type AgentResponse struct {
	ID               int64        `json:"ID"`
	Name             string       `json:"name"`
	Instructions     string       `json:"instructions"`
	WelcomeMessage   string       `json:"welcome_message"`
	SuggestedPrompts []string     `json:"suggested_prompts"`
	Files            []FileDetail `json:"files"`
	Status           string       `json:"status"`
	EmbeddingsStatus string       `json:"embeddings_status"`
	CreatedOn        int64        `json:"created_on"`
	UpdatedOn        int64        `json:"updated_on"`
}
