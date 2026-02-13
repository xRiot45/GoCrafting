package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/scaffold"
)

var generateCmd = &cobra.Command{
	Use:     "generate [schematic] [name]",
	Aliases: []string{"g"},
	Short:   "Generate boilerplate code based on schematic",
	Args:    cobra.MinimumNArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		schematic := args[0]
		name := args[1]

		meta, err := core.LoadMetadata()
		if err != nil {
			handleError(fmt.Errorf("gocrafting-cli.json not found. Are you in the root of the project?"))
		}

		fmt.Printf("🛠  Scaffolding %s '%s' for %s...\n", schematic, name, meta.SelectedFramework)
		start := time.Now()

		if err := scaffold.Run(meta, schematic, name); err != nil {
			handleError(err)
		}

		fmt.Printf("✅ Done in %s\n", time.Since(start))
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
