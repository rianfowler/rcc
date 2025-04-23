package test

import (
	"context"
	"testing"

	"dagger.io/dagger"
)

func TestDaggerBuild(t *testing.T) {
	ctx := context.Background()

	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(t.Log))
	if err != nil {
		t.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// Get the source directory
	src := client.Host().Directory(".")

	// Create a container with Go
	golang := client.Container().From("golang:1.24")

	// Mount the source directory
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

	// Run go build
	_, err = golang.WithExec([]string{"go", "build", "-o", "testapp", "test/testdata/go-app/main.go"}).Stdout(ctx)
	if err != nil {
		t.Fatalf("Failed to build Go application: %v", err)
	}

	// Verify the binary was created
	output, err := golang.WithExec([]string{"ls", "-l", "testapp"}).Stdout(ctx)
	if err != nil {
		t.Fatalf("Failed to verify binary: %v", err)
	}
	t.Logf("Build output:\n%s", output)
}
