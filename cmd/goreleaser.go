package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newGoreleaserBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go [tag]",
		Short: "Run goreleaser to create a release for the given tag",
		Long: `This command runs the goreleaser CLI to build and release your Go application.
It sets the GORELEASER_CURRENT_TAG environment variable and runs goreleaser in the specified directory.`,
		Args: cobra.ExactArgs(1),
		RunE: runGoreleaserBuild,
	}

	// Add path flag
	cmd.Flags().String("path", ".", "Path to the Go application to build")

	return cmd
}

func runGoreleaserBuild(cmd *cobra.Command, args []string) error {
	tag := args[0]
	sourcePath, _ := cmd.Flags().GetString("path")

	// Get absolute path of source directory
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Change to the source directory
	if err := os.Chdir(absPath); err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	// Set environment variable for the tag
	if err := os.Setenv("GORELEASER_CURRENT_TAG", tag); err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}

	// Run goreleaser
	goreleaserCmd := exec.Command("goreleaser", "release", "--snapshot", "--skip", "docker,homebrew", "--verbose")
	goreleaserCmd.Stdout = os.Stdout
	goreleaserCmd.Stderr = os.Stderr

	if err := goreleaserCmd.Run(); err != nil {
		return fmt.Errorf("failed to run goreleaser: %v", err)
	}

	return nil
}
