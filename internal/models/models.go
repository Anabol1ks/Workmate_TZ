package models

import "github.com/google/uuid"

type TaskStatus string

const (
	StatusCreated TaskStatus = "created"
	StatusActive  TaskStatus = "active"
	StatusDone    TaskStatus = "done"
	StatusError   TaskStatus = "error"
)

type Task struct {
	ID         uuid.UUID   `json:"id"`
	Status     TaskStatus  `json:"status"`
	Files      []FileInfo  `json:"files"`
	ArchiveURL string      `json:"archive_url,omitempty"`
	Errors     []FileError `json:"errors,omitempty"`
}

type FileInfo struct {
	URL        string `json:"url"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Downloader bool   `json:"downloader"`
}

type FileError struct {
	URL    string `json:"url"`
	Reason string `json:"reason"`
}
