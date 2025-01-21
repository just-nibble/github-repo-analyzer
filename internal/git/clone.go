// internal/git/clone.go
package git

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/just-nibble/github-repo-analyzer/internal/models"
)

// Service handles Git operations
type Service struct {
	tempDir string
}

// NewService creates a new Git service
func NewService() (*Service, error) {
	tempDir, err := os.MkdirTemp("", "repo-analysis-*")
	if err != nil {
		return nil, models.NewAnalysisError("failed to create temporary directory", err)
	}
	return &Service{tempDir: tempDir}, nil
}

// CloneRepo clones a repository and returns the path
func (s *Service) CloneRepo(url string) error {
	_, err := git.PlainClone(s.tempDir, false, &git.CloneOptions{
		URL:               url,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return models.NewAnalysisError("failed to clone repository", err)
	}
	return nil
}

// HasSubmodules checks if the repository has submodules
func (s *Service) HasSubmodules() bool {
	gitModules := filepath.Join(s.tempDir, ".gitmodules")
	_, err := os.Stat(gitModules)
	return err == nil
}

// GetRepoPath returns the path to the cloned repository
func (s *Service) GetRepoPath() string {
	return s.tempDir
}

// Cleanup removes the temporary directory
func (s *Service) Cleanup() error {
	if err := os.RemoveAll(s.tempDir); err != nil {
		return models.NewAnalysisError("failed to cleanup temporary directory", err)
	}
	return nil
}
