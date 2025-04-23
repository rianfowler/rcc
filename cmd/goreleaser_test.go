package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoreleaserBuild(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "goreleaser-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Copy test Go app to temp directory
	testAppPath := filepath.Join("test", "testdata", "go-app")
	err = copyDir(testAppPath, tempDir)
	assert.NoError(t, err)

	// Create a command with the test directory
	cmd := newGoreleaserBuildCmd()
	cmd.SetArgs([]string{"v1.0.0", "--path", tempDir})

	// Run the command
	err = cmd.Execute()
	assert.NoError(t, err)

	// Verify that the dist directory was created
	distPath := filepath.Join(tempDir, "dist")
	_, err = os.Stat(distPath)
	assert.NoError(t, err)

	// Verify that the config.yaml was created
	configPath := filepath.Join(distPath, "config.yaml")
	_, err = os.Stat(configPath)
	assert.NoError(t, err)
}

// Helper function to copy a directory
func copyDir(src, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Read source file
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			// Write destination file
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
