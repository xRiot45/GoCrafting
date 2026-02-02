// Package runner provides functions for executing commands in a project directory.
package runner

import (
	"fmt"
	"os"
	"os/exec"
)

// RunGoModTidy run the 'go mod tidy' command in the target directory
func RunGoModTidy(projectPath string) error {
	fmt.Println("📦 Downloading dependencies...")

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}
	return nil
}

// RunGoFmt run the 'go fmt' command in the target directory
func RunGoFmt(projectPath string) error {
	fmt.Println("🧹 Formatting code...")

	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go fmt: %w", err)
	}
	return nil
}

// RunGitInit runs the 'git init' command in the target directory
func RunGitInit(projectPath string) error {
	fmt.Println("create git repository...")

	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	return cmd.Run()
}

// GoGet installs packages using the 'go get' command
func GoGet(projectPath string, packages ...string) error {
	if len(packages) == 0 {
		return nil
	}

	fmt.Printf("⬇️  Installing packages: %v...\n", packages)

	args := append([]string{"get"}, packages...)

	// #nosec G204 -- Arguments are controlled internally, not direct user input
	cmd := exec.Command("go", args...)
	cmd.Dir = projectPath
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install packages: %w", err)
	}
	return nil
}
