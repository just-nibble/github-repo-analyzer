// internal/analyzer/analyzer_test.go
package analyzer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// MockGitService implements GitService for testing
type MockGitService struct {
	cloneErr      error
	hasSubmodules bool
	repoPath      string
	cleanupErr    error
}

func (m *MockGitService) CloneRepo(url string) error { return m.cloneErr }
func (m *MockGitService) HasSubmodules() bool        { return m.hasSubmodules }
func (m *MockGitService) GetRepoPath() string        { return m.repoPath }
func (m *MockGitService) Cleanup() error             { return m.cleanupErr }

func setupTestRepo(t *testing.T) (string, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "analyzer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create test file structure
	testFiles := map[string]int64{
		"file1.txt":           1024 * 1024,     // 1MB
		"file2.txt":           1024 * 1024 * 2, // 2MB
		"dir1/file3.txt":      1024 * 512,      // 0.5MB
		"dir1/dir2/file4.txt": 1024 * 1024 * 3, // 3MB
	}

	for path, size := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		data := make([]byte, size)
		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestAnalyzeRepo(t *testing.T) {
	tempDir, cleanup := setupTestRepo(t)
	defer cleanup()

	tests := []struct {
		name          string
		gitService    GitService
		expectError   bool
		expectedSize  float64
		expectedFiles int
	}{
		{
			name: "Successful Analysis",
			gitService: &MockGitService{
				repoPath:      tempDir,
				hasSubmodules: true,
			},
			expectError:   false,
			expectedSize:  6.5, // 1MB + 2MB + 0.5MB + 3MB
			expectedFiles: 4,
		},
		{
			name: "Clone Error",
			gitService: &MockGitService{
				cloneErr: errors.New("clone failed"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.gitService)
			analysis, err := service.AnalyzeRepo("test-url")

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify total size
			if analysis.Size != tt.expectedSize {
				t.Errorf("Got size %f MB, want %f MB", analysis.Size, tt.expectedSize)
			}

			// Count total files
			totalFiles := 0
			for _, folder := range analysis.Folders {
				totalFiles += len(folder.Files)
			}

			if totalFiles != tt.expectedFiles {
				t.Errorf("Got %d files, want %d files", totalFiles, tt.expectedFiles)
			}

			// Verify submodules detection
			if analysis.HasSubmodules != tt.gitService.HasSubmodules() {
				t.Errorf("Got HasSubmodules %v, want %v",
					analysis.HasSubmodules,
					tt.gitService.HasSubmodules())
			}
		})
	}
}

func TestAnalyzeDirectory(t *testing.T) {
	tempDir, cleanup := setupTestRepo(t)
	defer cleanup()

	service := NewService(&MockGitService{})
	folder, err := service.analyzeDirectory(filepath.Join(tempDir, "dir1"), tempDir)
	if err != nil {
		t.Fatalf("analyzeDirectory failed: %v", err)
	}

	// Verify folder name
	expectedName := "dir1"
	if folder.Name != expectedName {
		t.Errorf("Got folder name %s, want %s", folder.Name, expectedName)
	}

	// Verify files
	if len(folder.Files) != 1 {
		t.Errorf("Got %d files, want 1", len(folder.Files))
	}

	// Verify file size
	expectedSize := 0.5 // 512KB = 0.5MB
	if folder.Files[0].Size != expectedSize {
		t.Errorf("Got file size %f MB, want %f MB", folder.Files[0].Size, expectedSize)
	}
}
