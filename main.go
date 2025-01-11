// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// RepoAnalysis represents the JSON output structure
type RepoAnalysis struct {
	CloneURL      string   `json:"clone_url"`
	Size          float64  `json:"size"`       // Total size in MB
	SizeHuman     string   `json:"size_human"` // Human readable size
	Folders       []Folder `json:"folders"`
	HasSubmodules bool     `json:"has_submodules"`
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

// bytesToHumanReadable converts bytes to human readable format
func bytesToHumanReadable(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	size := float64(bytes) / float64(div)
	return fmt.Sprintf("%.2f %cB", size, "KMGTPE"[exp])
}

// bytesToMB converts bytes to megabytes
func bytesToMB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// analyzeDirectory traverses a directory and returns folder information
func analyzeDirectory(path string, baseDir string) (Folder, error) {
	folder := Folder{
		Name:  strings.TrimPrefix(path, baseDir+"/"),
		Files: []File{},
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return folder, fmt.Errorf("error reading directory %s: %v", path, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return folder, fmt.Errorf("error getting file info for %s: %v", entry.Name(), err)
			}

			fileSize := info.Size()
			sizeMB := bytesToMB(fileSize)

			file := File{
				Name:      entry.Name(),
				Size:      float64(int(sizeMB*100)) / 100, // Round to 2 decimal places
				SizeHuman: bytesToHumanReadable(fileSize),
			}
			folder.Files = append(folder.Files, file)
		}
	}

	return folder, nil
}

// hasSubmodules checks if the repository has any submodules
func hasSubmodules(repoPath string) bool {
	gitModules := filepath.Join(repoPath, ".gitmodules")
	_, err := os.Stat(gitModules)
	return err == nil
}

// analyzeRepo analyzes the cloned repository
func analyzeRepo(repoPath string) (RepoAnalysis, error) {
	analysis := RepoAnalysis{
		Folders: []Folder{},
	}

	var totalSize int64
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		if info.IsDir() && path != repoPath {
			folder, err := analyzeDirectory(path, repoPath)
			if err != nil {
				return err
			}
			analysis.Folders = append(analysis.Folders, folder)
		}
		return nil
	})

	if err != nil {
		return analysis, fmt.Errorf("error walking repository: %v", err)
	}

	analysis.Size = float64(int(bytesToMB(totalSize)*100)) / 100 // Round to 2 decimal places
	analysis.SizeHuman = bytesToHumanReadable(totalSize)
	analysis.HasSubmodules = hasSubmodules(repoPath)

	return analysis, nil
}

// cloneOptions creates git clone options with submodule support
func cloneOptions(url string) *git.CloneOptions {
	return &git.CloneOptions{
		URL:               url,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <github-repo-url>")
	}

	repoURL := os.Args[1]
	analysis := RepoAnalysis{
		CloneURL: repoURL,
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "repo-analysis-*")
	if err != nil {
		log.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Clone repository with submodules
	log.Printf("Cloning repository %s (including submodules)...", repoURL)
	_, err = git.PlainClone(tempDir, false, cloneOptions(repoURL))
	if err != nil {
		log.Fatalf("Error cloning repository: %v", err)
	}

	// Analyze repository
	repoAnalysis, err := analyzeRepo(tempDir)
	if err != nil {
		log.Fatalf("Error analyzing repository: %v", err)
	}

	analysis.Size = repoAnalysis.Size
	analysis.SizeHuman = repoAnalysis.SizeHuman
	analysis.Folders = repoAnalysis.Folders
	analysis.HasSubmodules = repoAnalysis.HasSubmodules

	// Output JSON
	output, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}

	fmt.Println(string(output))
}
