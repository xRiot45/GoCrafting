// Package cmd defines the entry point and root command for the CLI.
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// --- 1. STYLING CONSTANTS ---

var (
	colorCyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)  // Main Title
	colorGreen   = lipgloss.NewStyle().Foreground(lipgloss.Color("76")).Bold(true)  // Key Commands
	colorMagenta = lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Bold(true) // Categories
	colorWhite   = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))            // Standard Text
	colorGray    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // Descriptions
	colorDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("237"))            // Brackets/Decorations
	colorHeader  = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("62")).Bold(true).Padding(0, 1)
)

// --- 2. DATA STRUCTURES ---

type helpEntry struct {
	Name, Desc string
}

type schematicEntry struct {
	Name, Alias, Desc string
}

// --- 3. ROOT COMMAND DEFINITION ---

var rootCmd = &cobra.Command{
	Use:   "gocrafting",
	Short: "The Enterprise-Grade Go Architecture Generator",
	Run: func(cmd *cobra.Command, _ []string) {
		renderCustomHelp(cmd)
	},
}

// --- 4. INITIALIZATION ---

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		renderCustomHelp(cmd)
	})

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleError(err)
	}
}

// --- 6. CUSTOM HELP RENDERER (THE BRAIN) ---

func renderCustomHelp(_ *cobra.Command) {
	// Header Logo & Version
	fmt.Println()
	fmt.Println(colorCyan.Render("  ⚡ GO CRAFTING CLI ") + colorGray.Render(" v1.0.0"))
	fmt.Println(colorGray.Render("  ──────────────────────────"))
	fmt.Println("  The Enterprise-Grade Go Architecture Generator.")
	fmt.Println("  Manage your entire development lifecycle with precision.")
	fmt.Println()

	// Usage Section
	fmt.Println(colorHeader.Render("USAGE"))
	fmt.Println("  $ gocrafting <command> [subcommand] [flags]")
	fmt.Println()

	// SECTION 1: CORE COMMANDS
	printSection("CORE COMMANDS", []helpEntry{
		{"new (n)", "Create a new production-ready Go project."},
		{"info (i)", "Display project context (Framework, DB, Addons)."},
		{"update (u)", "Self-update the CLI to the latest version."},
		{"doctor", "Check your environment health (Go, Docker, Make)."},
	})

	// SECTION 2: DEV & OPS
	printSection("DEVELOPMENT & OPS", []helpEntry{
		{"run (r)", "Start dev server with hot-reload (Auto-configures Air)."},
		{"build (b)", "Compile optimized binary (Injects version/git-hash)."},
		{"test (t)", "Run unit tests with rich reporting coverage."},
		{"docker (d)", "Generate/Update Dockerfile & Compose configs."},
	})

	// SECTION 3: DATABASE
	printSection("DATABASE MANAGEMENT", []helpEntry{
		{"migrate", "Manage database migrations (create|up|down|version)."},
		{"seed", "Run database seeders to populate initial data."},
	})

	// SECTION 4: GENERATORS (SCHEMATICS)
	printSchematicsList()

	// SECTION 5: EXTENSIONS
	printSection("EXTENSION", []helpEntry{
		{"add <lib>", "Install a library & generate setup code (e.g. 'add redis')."},
	})

	// SECTION 6: FLAGS
	printSection("FLAGS", []helpEntry{
		{"-h, --help", "Show this help message."},
		{"-v, --version", "Show CLI version."},
		{"--verbose", "Enable detailed logging."},
	})

	// EXAMPLES
	fmt.Println(colorCyan.Render("  EXAMPLES"))
	fmt.Println("    $ gocrafting new my-saas")
	fmt.Println("    $ gocrafting g resource products")
	fmt.Println("    $ gocrafting migrate up")
	fmt.Println("    $ gocrafting add redis")
	fmt.Println()

	// FOOTER
	fmt.Println(colorGray.Render("  Need help? https://github.com/xRiot45/gocrafting/issues"))
	fmt.Println()
}

// --- 7. HELPER FUNCTIONS ---

// printSection prints a list of commands using TabWriter.
func printSection(title string, entries []helpEntry) {
	fmt.Println(colorCyan.Render("  " + title))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	for _, e := range entries {
		_, _ = fmt.Fprintf(w, "    %s\t%s\n", colorWhite.Render(e.Name), colorGray.Render(e.Desc))
	}
	_ = w.Flush()
	fmt.Println()
}

// printSchematicsList prints generators with grouping.
func printSchematicsList() {
	fmt.Println(colorCyan.Render("  GENERATORS (Schematics)"))
	fmt.Println("    Generate boilerplate code based on your chosen architecture.")
	fmt.Println("    Usage: gocrafting generate <schematic> [name]")
	fmt.Println()

	printGroup := func(groupTitle string, items []schematicEntry) {
		fmt.Println(colorMagenta.Render("    " + groupTitle))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		for _, item := range items {
			aliasFmt := fmt.Sprintf("%s%s%s", colorDim.Render("["), item.Alias, colorDim.Render("]"))
			_, _ = fmt.Fprintf(w, "      %s\t%s\t%s\n",
				colorGreen.Render(item.Name),
				aliasFmt,
				colorGray.Render(item.Desc),
			)
		}

		_ = w.Flush()
		fmt.Println()
	}

	// Group 1: Core Architecture
	printGroup("CORE ARCHITECTURE", []schematicEntry{
		{"resource", "res", "Full CRUD (Handler + Service + Repo + DTO)"},
		{"handler", "h", "HTTP Handler (Fiber/Gin/Echo context aware)"},
		{"service", "s", "Business Logic Layer (Interface & Impl)"},
	})

	// Group 2: Data Layer
	printGroup("DATA LAYER", []schematicEntry{
		{"repository", "repo", "Data Access Layer (SQL/Gorm implementation)"},
		{"model", "m", "Database Entity & DTO structs"},
		{"migration", "mig", "Database schema migration file"},
	})

	// Group 3: System & Config
	printGroup("SYSTEM & CONFIG", []schematicEntry{
		{"middleware", "mid", "HTTP Middleware template"},
		{"config", "conf", "Environment configuration loader"},
		{"docker", "d", "Dockerfile & Docker Compose setup"},
		{"proto", "grpc", "Protobuf definition & gRPC stub generation"},
		{"cron", "job", "Scheduled task/Cron job template"},
	})
}

// handleError handles execution errors with a styled red box.
func handleError(err error) {
	var (
		errBoxStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FF0000")).
				Padding(0, 1).
				MarginTop(1)

		errTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Bold(true)

		errBodyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))
	)

	errorMsg := fmt.Sprintf("%s\n%s",
		errTitleStyle.Render("❌ EXECUTION FAILED"),
		errBodyStyle.Render(err.Error()),
	)

	// Print to Stderr so it is separated from normal output
	fmt.Fprintln(os.Stderr, errBoxStyle.Render(errorMsg))
	os.Exit(1)
}
