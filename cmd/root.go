// Package cmd handles the command-line interface logic using Cobra.
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xRiot45/gocrafting/internal/ui"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "gocrafting",
	Short: "The Artisan's Scaffolder for Go applications",
	Long: `GoCrafting is a specialized CLI tool designed for developers who treat coding as an art. 
It bridges the gap between raw ideas and production-ready code by scaffolding 
consistent, idiomatic, and enterprise-ready project structures.`,
}

// initCmd represents the 'init' command used to start a new project.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start a new project canvas",
	// We use _ for unused parameters to satisfy the linter.
	Run: func(_ *cobra.Command, _ []string) {
		// Initialize and run the Bubble Tea TUI program.
		p := tea.NewProgram(ui.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
func Execute() {
	// Register the 'init' sub-command.
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
