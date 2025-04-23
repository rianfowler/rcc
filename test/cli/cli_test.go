package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rianfowler/rcc/test/testutil"
)

func TestBuildPackCommand(t *testing.T) {
	ctx, _ := testutil.TestContext(t)
	testDataPath := testutil.GetTestDataPath(t)

	// Build the Node.js app using the CLI
	cmd := exec.CommandContext(ctx, "go", "run", "../..", "build", filepath.Join(testDataPath, "node-app"), "--name", "test-node-app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run build command: %v", err)
	}

	// Verify the image was built
	dockerCmd := exec.CommandContext(ctx, "docker", "image", "inspect", "test-node-app")
	if err := dockerCmd.Run(); err != nil {
		t.Fatalf("docker image not found: %v", err)
	}

	// Clean up
	cleanupCmd := exec.CommandContext(ctx, "docker", "rmi", "test-node-app")
	cleanupCmd.Run() // Ignore error as image might not exist
}

func TestBuildPackExplicitCommand(t *testing.T) {
	ctx, _ := testutil.TestContext(t)
	testDataPath := testutil.GetTestDataPath(t)

	// Build the Node.js app using the explicit pack subcommand
	cmd := exec.CommandContext(ctx, "go", "run", "../..", "build", "pack", filepath.Join(testDataPath, "node-app"), "--name", "test-node-app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run build pack command: %v", err)
	}

	// Verify the image was built
	dockerCmd := exec.CommandContext(ctx, "docker", "image", "inspect", "test-node-app")
	if err := dockerCmd.Run(); err != nil {
		t.Fatalf("docker image not found: %v", err)
	}

	// Clean up
	cleanupCmd := exec.CommandContext(ctx, "docker", "rmi", "test-node-app")
	cleanupCmd.Run() // Ignore error as image might not exist
}

func TestBuildGoreleaserCommand(t *testing.T) {
	ctx, _ := testutil.TestContext(t)
	testDataPath := testutil.GetTestDataPath(t)

	// Run goreleaser using the CLI with our test Go app
	cmd := exec.CommandContext(ctx, "go", "run", "../..", "build", "go", "v0.0.0-test", "--path", filepath.Join(testDataPath, "go-app"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run goreleaser command: %v", err)
	}

	// Verify dist directory exists
	distPath := "dist"
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		t.Fatalf("dist directory not found: %v", err)
	}

	// Check for basic expected files
	expectedFiles := []string{
		"checksums.txt",
		"config.yaml",
	}

	files, err := os.ReadDir(distPath)
	if err != nil {
		t.Fatalf("failed to read dist directory: %v", err)
	}

	for _, expected := range expectedFiles {
		found := false
		for _, file := range files {
			if file.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file %s not found in dist directory", expected)
		}
	}

	// Clean up
	os.RemoveAll(distPath)
}
