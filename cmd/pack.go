package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newPackBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pack [path]",
		Short: "Build an application using Cloud Native Buildpacks",
		Long: `Build an application using Cloud Native Buildpacks.
This command will create a container image from your application source code.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			// Validate path exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", path)
			}
			return runPackBuild(cmd, args)
		},
	}

	// Add pack-specific flags
	cmd.Flags().String("name", "", "Name for the output image (required)")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runPackBuild(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]

	// Get absolute path of source directory
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Initialize Docker client with specific API version
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.48"),
	)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer cli.Close()

	// Pull the pack image if not available
	imageName := "buildpacksio/pack:latest"
	_, err = cli.ImageInspect(context.Background(), imageName)
	if err != nil {
		log.Printf("Pulling image %s...", imageName)
		out, err := cli.ImagePull(context.Background(), imageName, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image: %v", err)
		}
		defer out.Close()
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}

	// Get flags
	name, _ := cmd.Flags().GetString("name")
	builder := viper.GetString("builder")
	envVars := viper.GetStringSlice("env")

	// Prepare command arguments
	args = []string{"build", name}
	args = append(args, "--path", "/workspace")
	args = append(args, "--builder", builder)
	args = append(args, "--creation-time", "now")

	// Add environment variables
	for _, env := range envVars {
		args = append(args, "--env", env)
	}

	// Create container config
	config := &container.Config{
		Image: imageName,
		Cmd:   args,
		User:  "root", // Run as root to ensure access to Docker socket
	}

	// Create host config with volume mount
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/workspace", absPath),
			"/var/run/docker.sock:/var/run/docker.sock",
		},
		// Ensure the container has access to the Docker socket
		SecurityOpt: []string{"label:disable"},
	}

	// Create the container
	resp, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
	}

	// Get the container logs
	out, err := cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return fmt.Errorf("failed to get container logs: %v", err)
	}

	// Copy the logs to stdout and stderr
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// Remove the container
	if err := cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{}); err != nil {
		log.Printf("Warning: Failed to remove container: %v", err)
	}

	return nil
}
