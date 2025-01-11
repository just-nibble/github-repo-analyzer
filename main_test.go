// main_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBytesToMB(t *testing.T) {
	tests := []struct {
		input    int64
		expected float64
	}{
		{1024 * 1024, 1.0},
		{1024 * 1024 * 2, 2.0},
		{1024 * 1024 * 1.5, 1.5},
	}

	for _, test := range tests {
		result := bytesToMB(test.input)
		if result != test.expected {
			t.Errorf("bytesToMB(%d) = %f; want %f", test.input, result, test.expected)
		}
	}
}

func TestBytesToHumanReadable(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
		{1024 * 1024 * 1024 * 1024, "1.00 TB"},
	}

	for _, test := range tests {
		result := bytesToHumanReadable(test.input)
		if result != test.expected {
			t.Errorf("bytesToHumanReadable(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestAnalyzeDirectory(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "test-repo-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]int64{
		"file1.txt": 1024 * 1024,     // 1MB
		"file2.txt": 1024 * 1024 * 2, // 2MB
	}

	for name, size := range testFiles {
		path := filepath.Join(tempDir, name)
		data := make([]byte, size)
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Test directory analysis
	folder, err := analyzeDirectory(tempDir, filepath.Dir(tempDir))
	if err != nil {
		t.Fatalf("analyzeDirectory failed: %v", err)
	}

	// Verify results
	if len(folder.Files) != len(testFiles) {
		t.Errorf("Got %d files, want %d", len(folder.Files), len(testFiles))
	}

	for _, file := range folder.Files {
		expectedSize, exists := testFiles[file.Name]
		if !exists {
			t.Errorf("Unexpected file: %s", file.Name)
			continue
		}
		expectedMB := bytesToMB(expectedSize)
		if file.Size != expectedMB {
			t.Errorf("File %s: got size %f MB, want %f MB", file.Name, file.Size, expectedMB)
		}

		// Verify human readable size
		expectedHuman := bytesToHumanReadable(expectedSize)
		if file.SizeHuman != expectedHuman {
			t.Errorf("File %s: got human size %s, want %s", file.Name, file.SizeHuman, expectedHuman)
		}
	}
}

func TestHasSubmodules(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "test-repo-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initially should have no submodules
	if hasSubmodules(tempDir) {
		t.Error("New directory should not have submodules")
	}

	// Create .gitmodules file
	gitModulesPath := filepath.Join(tempDir, ".gitmodules")
	if err := os.WriteFile(gitModulesPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create .gitmodules file: %v", err)
	}

	// Now should detect submodules
	if !hasSubmodules(tempDir) {
		t.Error("Should detect .gitmodules file")
	}
}
