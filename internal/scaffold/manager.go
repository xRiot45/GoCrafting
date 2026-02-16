package scaffold

import (
	"fmt"

	"github.com/xRiot45/gocrafting/internal/core"
)

// Run adalah traffic controller
func Run(meta *core.ProjectMetadata, schematic, name string) error {

	switch schematic {

	// --- CORE ---
	// case "resource", "res":
	// 	return GenerateResource(meta, name)

	case "handler", "h":
		return GenerateHandler(meta, name)

	// case "service", "s":
	// 	return GenerateService(meta, name)

	// // --- DATA LAYER ---
	// case "repository", "repo":
	// 	return GenerateRepository(meta, name)

	// case "model", "m":
	// 	return GenerateModel(meta, name)

	// case "migration", "mig":
	// 	return GenerateMigration(meta, name)

	// // --- SYSTEM ---
	// case "docker", "d":
	// 	return GenerateDocker(meta) // Mungkin tidak butuh param 'name'

	// case "middleware", "mid":
	// 	return GenerateMiddleware(meta, name)

	// case "cron", "job":
	// 	return GenerateCronJob(meta, name)

	default:
		return fmt.Errorf("unknown schematic: '%s'. Run 'gocrafting --help' to see available generators", schematic)
	}
}
