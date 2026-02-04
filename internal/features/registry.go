// Package features provides a registry to manage different project scale providers.
package features

import (
	"fmt"

	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/features/small"
)

// GetProvider returns a FeatureProvider based on the given project scale.
func GetProvider(scale string) (core.FeatureProvider, error) {
	switch scale {
	case "Small":
		return small.NewProvider(), nil

	case "Medium":
		// return medium.NewProvider(), nil
		return nil, fmt.Errorf("feature medium is coming soon")

	case "Enterprise":
		return nil, fmt.Errorf("feature enterprise is coming soon")

	default:
		return nil, fmt.Errorf("unknown project scale: %s", scale)
	}
}
