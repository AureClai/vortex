package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the Vortex application into a Wasm module.",
	Long: `Compiles the Go source code into a WebAssembly module (app.wasm) and
copies the necessary wasm_exec.js file. This command should be run from
the root of a Vortex project.`,
	Run: runBuild,
}

// runBuild handles the logic for the 'vortex build' command.
func runBuild(cmd *cobra.Command, args []string) {
	// Check if we are in a vortex project
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		fmt.Println("❌ No main.go file found. Are you in a Vortex project directory?")
		os.Exit(1)
	}

	fmt.Println("Building Go code to WebAssembly...")

	// Set environment variables for the build command.
	buildCmd := exec.Command("go", "build", "-o", "app.wasm", ".")
	buildCmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	// Run the build command.
	if err := buildCmd.Run(); err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Build successful.")

	// Copy the wasm_exec.js file.
	if err := copyWasmExec(); err != nil {
		fmt.Printf("❌ Failed to copy wasm_exec.js: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Copied wasm_exec.js.")
	fmt.Println("\nBuild complete. You can now serve the directory using 'vortex dev'")
}

// copyWasmExec finds and copies the wasm_exec.js file.
func copyWasmExec() error {
	goRoot := os.Getenv("GOROOT")
	if goRoot == "" {
		return fmt.Errorf("GOROOT environment variable is not set")
	}

	srcPath := filepath.Join(goRoot, "lib", "wasm", "wasm_exec.js")
	destPath := "wasm_exec.js"

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("could not open source file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("could not create destination file %s: %w", destPath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("could not copy file contents: %w", err)
	}
	return nil
}
