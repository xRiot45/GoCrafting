package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xRiot45/gocrafting/internal/ui"
)

var newCmd = &cobra.Command{
	Use:     "new [project-name]",
	Aliases: []string{"n"},
	Short:   "Create a new production-ready Go project.",
	Long:    "Starts the interactive TUI to generate a clean architecture Go project.",
	Example: "  gocrafting new\n  gocrafting new my-service",
	Run: func(_ *cobra.Command, args []string) {

		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		if err := ui.Start(projectName); err != nil {
			handleError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
