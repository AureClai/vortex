package main

import "github.com/spf13/cobra"

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vortex",
	Short: "Vortex is a front-end web framework for Go and WebAssembly.",
	Long: `Vortex provides a CLI to initialize, build, and serve Go-based
front-end applications that compile to WebAssembly.`,
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(devCmd)
}
