// internal/git/clone_test.go
package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewService(t *testing.T) {
	// Test service creation
	service, err := NewService()
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}
	defer service.Cleanup()

	// Verify temp directory was created
	if _, err := os.Stat(service.tempDir); os.IsNotExist(err) {
		t.Error("Temporary directory was not created")
	}

	// Verify temp directory has correct prefix
	base := filepath.Base(service.tempDir)
	if len(base) < len("repo-analysis-") || base[:len("repo-analysis-")] != "repo-analysis-" {
		t.Error("Temporary directory has incorrect prefix")
	}
}

func TestCleanup(t *testing.T) {
	service, err := NewService()
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	tempDir := service.tempDir
	if err := service.Cleanup(); err != nil {
		t.Errorf("Cleanup() failed: %v", err)
	}

	// Verify directory was removed
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Temporary directory was not removed")
	}
}

func TestHasSubmodules(t *testing.T) {
	service, err := NewService()
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}
	defer service.Cleanup()

	// Test without .gitmodules
	if service.HasSubmodules() {
		t.Error("HasSubmodules() should return false when no .gitmodules exists")
	}

	// Create .gitmodules file
	gitModulesPath := filepath.Join(service.tempDir, ".gitmodules")
	if err := os.WriteFile(gitModulesPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test .gitmodules: %v", err)
	}

	// Test with .gitmodules
	if !service.HasSubmodules() {
		t.Error("HasSubmodules() should return true when .gitmodules exists")
	}
}

func TestGetRepoPath(t *testing.T) {
	service, err := NewService()
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}
	defer service.Cleanup()

	path := service.GetRepoPath()
	if path != service.tempDir {
		t.Errorf("GetRepoPath() = %s; want %s", path, service.tempDir)
	}
}

func TestCloneRepo(t *testing.T) {
	service, err := NewService()
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}
	defer service.Cleanup()

	// Test with invalid URL
	err = service.CloneRepo("invalid-url")
	if err == nil {
		t.Error("CloneRepo() should fail with invalid URL")
	}

	// Note: Testing with a real repository would make this test slow and dependent
	// on external services. In a real project, you might want to:
	// 1. Use a mock git service for unit tests
	// 2. Have separate integration tests for testing with real repositories
}
