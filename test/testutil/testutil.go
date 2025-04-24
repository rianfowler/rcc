package testutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
)

// TestContext creates a new Dagger client and context for testing
func TestContext(t *testing.T) (context.Context, *dagger.Client) {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Fatalf("failed to connect to Dagger: %v", err)
	}
	t.Cleanup(func() {
		client.Close()
	})
	return ctx, client
}

// GetTestDataPath returns the absolute path to the testdata directory
func GetTestDataPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	return filepath.Join(wd, "..", "testdata")
}

// AssertContainerOutput runs a container and asserts its output contains the expected string
func AssertContainerOutput(t *testing.T, container *dagger.Container, expected string) {
	output, err := container.Stdout(context.Background())
	if err != nil {
		t.Fatalf("failed to get container output: %v", err)
	}
	if output != expected {
		t.Errorf("expected output %q, got %q", expected, output)
	}
}

// CreateTempDir creates a temporary directory for test artifacts
func CreateTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "rcc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}
