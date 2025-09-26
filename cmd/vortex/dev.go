package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs a local development server for the Vortex application.",
	Long: `Starts a static file server on port 8080 to serve the application.
It is recommended to run 'vortex build' before starting the dev server.`,
	Run: runDev,
}

// runDev handles the logic for the 'vortex dev' command.
func runDev(cmd *cobra.Command, args []string) {
	// Check if the build artifacts exist
	if _, err := os.Stat("app.wasm"); os.IsNotExist(err) {
		fmt.Println("⚠️ app.wasm not found. Did you run 'vortex build' first?")
	}

	port := "8080"
	addr := ":" + port
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Printf("Starting server on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
