// cmd/analyzer/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/just-nibble/github-repo-analyzer/internal/analyzer"
	"github.com/just-nibble/github-repo-analyzer/internal/git"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: github-repo-analyzer <github-repo-url>")
	}

	repoURL := os.Args[1]

	// Initialize services
	gitService, err := git.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize git service: %v", err)
	}

	analyzerService := analyzer.NewService(gitService)

	// Analyze repository
	analysis, err := analyzerService.AnalyzeRepo(repoURL)
	if err != nil {
		log.Printf("Error analyzing repository: %v", err)
		analysis.Error = err.Error()
	}

	// Output JSON
	output, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}

	fmt.Println(string(output))
}
