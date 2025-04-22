package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {
	// Get the source directory from command line args
	if len(os.Args) < 2 {
		log.Fatal("Please provide the source directory as an argument")
	}
	sourceDir := os.Args[1]

	// Get absolute path of source directory
	absPath, err := filepath.Abs(sourceDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	// Pull the pack image if not available
	imageName := "buildpacksio/pack:latest"
	_, _, err = cli.ImageInspectWithRaw(context.Background(), imageName)
	if err != nil {
		log.Printf("Pulling image %s...", imageName)
		out, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
		if err != nil {
			log.Fatalf("Failed to pull image: %v", err)
		}
		io.Copy(os.Stdout, out)
		out.Close()
	}

	// Create container config
	config := &container.Config{
		Image: imageName,
		Cmd: []string{
			"build", "demo-node-app",
			"--path", "/workspace",
			"--builder", "paketobuildpacks/builder-jammy-base",
			"--env", "BP_NODE_VERSION=18",
			"--env", "BP_DISABLE_SBOM=true",
		},
	}

	// Create host config with volume mount
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/workspace", absPath),
			"/var/run/docker.sock:/var/run/docker.sock",
		},
	}

	// Create the container
	resp, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, "")
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("Error waiting for container: %v", err)
		}
	case <-statusCh:
	}

	// Get the container logs
	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Fatalf("Failed to get container logs: %v", err)
	}

	// Copy the logs to stdout and stderr
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// Remove the container
	if err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{}); err != nil {
		log.Printf("Warning: Failed to remove container: %v", err)
	}
}
