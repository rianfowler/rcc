package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "rcc",
	Short: "RCC - Release and Container CLI tool",
	Long: `RCC is a CLI tool for building and releasing applications using
various build tools like Cloud Native Buildpacks and GoReleaser.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands
	rootCmd.AddCommand(NewBuildCmd())
}

func NewBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [path]",
		Short: "Build an application using various build tools",
		Long: `Build command supports multiple build tools:
- Cloud Native Buildpacks (default)
- GoReleaser

Examples:
  # Build using Cloud Native Buildpacks (default)
  rcc build ./path/to/app
  rcc build pack ./path/to/app

  # Build using GoReleaser
  rcc build go v1.0.0`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is specified and a path is provided, use pack
			if len(args) == 1 {
				return runPackBuild(cmd, args)
			}
			return cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(newPackBuildCmd())
	cmd.AddCommand(newGoreleaserBuildCmd())

	// Add flags
	cmd.PersistentFlags().String("builder", "paketobuildpacks/builder-jammy-base", "Builder to use for Cloud Native Buildpacks")
	cmd.PersistentFlags().StringSlice("env", []string{}, "Environment variables to pass to the build")
	cmd.Flags().String("name", "", "Name for the output image (required)")
	cmd.MarkFlagRequired("name")

	// Bind flags to viper
	viper.BindPFlag("builder", cmd.PersistentFlags().Lookup("builder"))
	viper.BindPFlag("env", cmd.PersistentFlags().Lookup("env"))

	return cmd
}
