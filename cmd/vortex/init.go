package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initializes a new Vortex project.",
	Long: `Creates a new directory with the specified project name and populates it with
the basic structure and files needed to get started with a Vortex application.
Features advanced Vortex styling with precompilation for optimal performance.`,
	Args: cobra.ExactArgs(1), // Ensures exactly one argument (the project name) is passed
	Run:  runInit,
}

func isValidProjectName(projectName string) bool {
	// Check for path separators (both Unix and Windows)
	if strings.ContainsAny(projectName, "/\\") {
		return false
	}
	// check for relative path attempts
	if strings.Contains(projectName, "..") {
		return false
	}
	// Checck for empty or whitespace only names
	return strings.TrimSpace(projectName) != ""
}

// runInit is the function executed when the 'init' command is called.
func runInit(cmd *cobra.Command, args []string) {
	projectName := args[0]

	fmt.Printf("🚀 Initializing new Vortex project: %s\n", projectName)
	fmt.Printf("✨ Including advanced style precompilation for maximum performance\n")

	// Validate project name to avoid directory traversal attacks
	if !isValidProjectName(projectName) {
		fmt.Println("❌ Invalid project name")
		return
	}

	// Create project directory
	if err := os.Mkdir(projectName, 0755); err != nil {
		fmt.Printf("❌ Directory %s already exists\n", projectName)
		return
	}
	fmt.Printf("✅ Created project directory: %s\n", projectName)

	// Create the subdirectory
	appDir := filepath.Join(projectName, "layout")
	if err := os.Mkdir(appDir, 0755); err != nil {
		fmt.Printf("❌ Directory %s already exists\n", appDir)
		return
	}
	fmt.Printf("✅ Created app directory: %s\n", appDir)

	// Create styles directory for organized styling
	stylesDir := filepath.Join(projectName, "styles")
	if err := os.Mkdir(stylesDir, 0755); err != nil {
		fmt.Printf("❌ Directory %s already exists\n", stylesDir)
		os.RemoveAll(projectName) // Cleanup
		os.Exit(1)
	}
	fmt.Printf("✅ Created styles directory: %s\n", stylesDir)

	// Define files to create with their content
	filesToCreate := map[string]string{
		"main.go":          getTemplate("main.go.tmpl", projectName),
		"index.html":       indexHTMLTemplate,
		"go.mod":           getTemplate("go.mod.tmpl", projectName),
		"README.md":        getTemplate("README.md.tmpl", projectName),
		".gitignore":       gitignoreTemplate,
		"layout/app.go":    getTemplate("app.go.tmpl", projectName),
		"styles/app.go":    getTemplate("styles.go.tmpl", projectName), // ✅ NEW: Dedicated styles file
		"styles/common.go": getTemplate("common.go.tmpl", projectName), // ✅ NEW: Common style patterns
	}

	for fileName, content := range filesToCreate {
		filePath := filepath.Join(projectName, fileName)
		err := os.WriteFile(filePath, []byte(strings.TrimSpace(content)), 0644)
		if err != nil {
			fmt.Printf("❌ Error creating file %s: %v\n", fileName, err)
			// Cleanup: attempt to remove created directory
			os.RemoveAll(projectName)
			return
		}
		fmt.Printf("✅ Created %s\n", filePath)
	}

	fmt.Printf("\n🎉 Project '%s' created successfully!\n\n", projectName)
	fmt.Println("🚀 Advanced Features Included:")
	fmt.Println("  ⚡ Style precompilation for maximum performance")
	fmt.Println("  🎨 Common style patterns ready to use")
	fmt.Println("  📊 Built-in performance monitoring")
	fmt.Println("  🔧 Organized style architecture")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. cd %s\n", projectName)
	fmt.Println("  2. Build the application: 'vortex build'")
	fmt.Println("  3. Start the dev server: 'vortex dev'")
	fmt.Println("  4. Open http://localhost:8080 in your browser.")
	fmt.Println("  5. Edit the layout/app.go file to start building your application.")
	fmt.Println("  6. Check styles/ directory for advanced styling patterns.")
}

// --- Enhanced Templates ---

//go:embed templates/*.tmpl
var templatesFS embed.FS

type TemplateData struct {
	ProjectName string
}

func getTemplate(name, projectName string) string {
	content, err := templatesFS.ReadFile("templates/" + name)
	if err != nil {
		fmt.Printf("❌ Error reading template %s: %v\n", name, err)
		return ""
	}
	return processTemplate(string(content), TemplateData{ProjectName: projectName})
}

func processTemplate(content string, data TemplateData) string {
	tmpl := template.Must(template.New("").Parse(content))
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Printf("❌ Error executing template: %v\n", err)
		return ""
	}
	return buf.String()
}

// Enhanced gitignore
const gitignoreTemplate = `
# Compiled Wasm file
*.wasm
app.wasm

# JS glue file  
wasm_exec.js

# Build artifacts
dist/
build/

# Editor files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Logs
*.log

# Performance profiling
*.prof
*.pprof
`

const indexHTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vortex App - High Performance</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
	<style>
 		* {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #2196F3 0%, #FF9800 100%);
            min-height: 100vh;
            color: #2C3E50;
            line-height: 1.6;
        }
        
        /* Performance indicator */
        .perf-indicator {
            position: fixed;
            top: 10px;
            right: 10px;
            background: rgba(76, 175, 80, 0.9);
            color: white;
            padding: 5px 10px;
            border-radius: 15px;
            font-size: 12px;
            font-weight: 600;
        }
	</style>
	</head>
<body>
    <div id="app"></div>
    <div class="perf-indicator">⚡ High Performance Mode</div>

    <script src="wasm_exec.js"></script>
    <script>
        if (!WebAssembly.instantiateStreaming) {
            WebAssembly.instantiateStreaming = async (resp, importObject) => {
                const source = await (await resp).arrayBuffer();
                return await WebAssembly.instantiate(source, importObject);
            };
        }

        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        }).catch((err) => {
            console.error("Wasm instantiation failed:", err);
            document.getElementById('app').innerHTML = 
                '<div class="error"><h2>⚠️ Failed to Load</h2><p>Check the browser console for details.</p></div>';
        });
    </script>
</body>
</html>
`
