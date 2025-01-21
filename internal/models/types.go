package models

import "fmt"

// RepoAnalysis represents the complete analysis of a repository
type RepoAnalysis struct {
	CloneURL      string   `json:"clone_url"`
	Size          float64  `json:"size"`       // Total size in MB
	SizeHuman     string   `json:"size_human"` // Human readable size
	Folders       []Folder `json:"folders"`
	HasSubmodules bool     `json:"has_submodules"`
	Error         string   `json:"error,omitempty"` // Error message if any
}

// Folder represents a directory in the repository
type Folder struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

// File represents a file in the repository
type File struct {
	Name      string  `json:"name"`
	Size      float64 `json:"size"`       // Size in MB
	SizeHuman string  `json:"size_human"` // Human readable size
}

// AnalysisError represents an error that occurred during analysis
type AnalysisError struct {
	Message string
	Cause   error
}

func (e *AnalysisError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewAnalysisError creates a new AnalysisError
func NewAnalysisError(message string, cause error) *AnalysisError {
	return &AnalysisError{
		Message: message,
		Cause:   cause,
	}
}
