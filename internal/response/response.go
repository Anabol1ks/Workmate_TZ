package response

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskStatusResponse struct {
	ID         string         `json:"id"`
	Status     string         `json:"status"`
	Files      []FileInfoDTO  `json:"files"`
	ArchiveURL string         `json:"archive_url,omitempty"`
	Errors     []FileErrorDTO `json:"errors,omitempty"`
}

type FileInfoDTO struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type FileErrorDTO struct {
	URL    string `json:"url"`
	Reason string `json:"reason"`
}
