// Package shell provides functions for executing commands in a project directory.
package shell

import (
	"fmt"
	"os/exec"
)

// GoGet installs the specified Go packages using 'go get'.
func GoGet(projectPath string, packages ...string) error {
	if len(packages) == 0 {
		return nil
	}

	args := append([]string{"get"}, packages...)

	// #nosec G204 -- Arguments are controlled internally by the generator, safe from injection.
	cmd := exec.Command("go", args...)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install packages: %w", err)
	}
	return nil
}

// RunGoModTidy executes 'go mod tidy' to clean up dependencies.
func RunGoModTidy(projectPath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}
	return nil
}

// RunGoFmt executes 'go fmt ./...' to format the project code.
func RunGoFmt(projectPath string) error {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go fmt: %w", err)
	}
	return nil
}

// RunGitInit initializes a new git repository in the project path.
func RunGitInit(projectPath string) error {
	fmt.Println("create git repository...")

	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	return cmd.Run()
}
