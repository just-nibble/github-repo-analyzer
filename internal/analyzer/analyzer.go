// internal/analyzer/analyzer.go
package analyzer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/just-nibble/github-repo-analyzer/internal/models"
	"github.com/just-nibble/github-repo-analyzer/pkg/utils"
)

// Service handles repository analysis
type Service struct {
	gitService GitService
}

// GitService interface for git operations
type GitService interface {
	CloneRepo(url string) error
	HasSubmodules() bool
	GetRepoPath() string
	Cleanup() error
}

// NewService creates a new analyzer service
func NewService(gitService GitService) *Service {
	return &Service{
		gitService: gitService,
	}
}

// analyzeDirectory analyzes a single directory and all its files
func (s *Service) analyzeDirectory(path, baseDir string) (models.Folder, error) {
	folder := models.Folder{
		Name:  strings.TrimPrefix(path, baseDir+string(os.PathSeparator)),
		Files: []models.File{},
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return folder, models.NewAnalysisError("failed to read directory", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return folder, models.NewAnalysisError("failed to get file info", err)
			}

			fileSize := info.Size()
			sizeMB := utils.RoundToTwoDecimals(utils.BytesToMB(fileSize))

			file := models.File{
				Name:      entry.Name(),
				Size:      sizeMB,
				SizeHuman: utils.BytesToHumanReadable(fileSize),
			}
			folder.Files = append(folder.Files, file)
		}
	}

	return folder, nil
}

// processDirectory recursively processes directories and collects all folders
func (s *Service) processDirectory(path string, baseDir string) ([]models.Folder, error) {
	var folders []models.Folder

	err := filepath.WalkDir(path, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return models.NewAnalysisError("failed to walk path", err)
		}

		if d.IsDir() {
			if currentPath != path { // Skip the root directory
				folder, err := s.analyzeDirectory(currentPath, baseDir)
				if err != nil {
					return err
				}
				folders = append(folders, folder)
			}
		}
		return nil
	})

	return folders, err
}

// calculateTotalSize calculates the total size of all files in the directory
func (s *Service) calculateTotalSize(path string) (int64, error) {
	var totalSize int64
	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return models.NewAnalysisError("failed to walk path", err)
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return models.NewAnalysisError("failed to get file info", err)
			}
			totalSize += info.Size()
		}
		return nil
	})
	return totalSize, err
}

// AnalyzeRepo analyzes a repository from URL
func (s *Service) AnalyzeRepo(url string) (models.RepoAnalysis, error) {
	analysis := models.RepoAnalysis{
		CloneURL: url,
		Folders:  []models.Folder{},
	}

	// Clone repository
	if err := s.gitService.CloneRepo(url); err != nil {
		return analysis, err
	}
	defer s.gitService.Cleanup()

	repoPath := s.gitService.GetRepoPath()

	// Calculate total size
	totalSize, err := s.calculateTotalSize(repoPath)
	if err != nil {
		return analysis, err
	}

	// Process all directories
	folders, err := s.processDirectory(repoPath, repoPath)
	if err != nil {
		return analysis, err
	}

	analysis.Size = utils.RoundToTwoDecimals(utils.BytesToMB(totalSize))
	analysis.SizeHuman = utils.BytesToHumanReadable(totalSize)
	analysis.HasSubmodules = s.gitService.HasSubmodules()
	analysis.Folders = folders

	return analysis, nil
}
