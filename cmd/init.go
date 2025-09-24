package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vortex",
	Short: "Vortex is a front-end web framework for Go and WebAssembly.",
	Long: `Vortex provides a CLI to initialize, build, and serve Go-based
front-end applications that compile to WebAssembly.`,
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initializes a new Vortex project.",
	Long: `Creates a new directory with the specified project name and populates it with
the basic structure and files needed to get started with a Vortex application.`,
	Args: cobra.ExactArgs(1), // Ensures exactly one argument (the project name) is passed
	Run:  runInit,
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the Vortex application into a Wasm module.",
	Long: `Compiles the Go source code into a WebAssembly module (app.wasm) and
copies the necessary wasm_exec.js file. This command should be run from
the root of a Vortex project.`,
	Run: runBuild,
}

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs a local development server for the Vortex application.",
	Long: `Starts a static file server on port 8080 to serve the application.
It is recommended to run 'vortex build' before starting the dev server.`,
	Run: runDev,
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(devCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// runInit is the function executed when the 'init' command is called.
func runInit(cmd *cobra.Command, args []string) {
	projectName := args[0]

	fmt.Printf("🚀 Initializing new Vortex project: %s\n", projectName)

	// Create project directory
	if err := os.Mkdir(projectName, 0755); err != nil {
		fmt.Printf("❌ Error creating project directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Created project directory: %s\n", projectName)

	// Create the subdirectory
	appDir := filepath.Join(projectName, "layout")
	if err := os.Mkdir(appDir, 0755); err != nil {
		fmt.Printf("❌ Error creating app directory: %v\n", err)
		os.RemoveAll(projectName) // Cleanup
		os.Exit(1)
	}
	fmt.Printf("✅ Created app directory: %s\n", appDir)

	// Define files to create with their content
	filesToCreate := map[string]string{
		"main.go":       mainGoTemplate(projectName),
		"index.html":    indexHTMLTemplate,
		"go.mod":        goModTemplate(projectName),
		"README.md":     readmeTemplate(projectName),
		".gitignore":    gitignoreTemplate,
		"layout/app.go": appGoTemplate,
	}

	for fileName, content := range filesToCreate {
		filePath := filepath.Join(projectName, fileName)
		err := os.WriteFile(filePath, []byte(strings.TrimSpace(content)), 0644)
		if err != nil {
			fmt.Printf("❌ Error creating file %s: %v\n", fileName, err)
			// Cleanup: attempt to remove created directory
			os.RemoveAll(projectName)
			os.Exit(1)
		}
		fmt.Printf("✅ Created %s\n", filePath)
	}

	fmt.Printf("\n🎉 Project '%s' created successfully!\n\n", projectName)
	fmt.Println("Next steps:")
	fmt.Printf("  1. cd %s\n", projectName)
	fmt.Println("  2. Build the application: 'vortex build'")
	fmt.Println("  3. Start the dev server: 'vortex dev'")
	fmt.Println("  4. Open http://localhost:8080 in your browser.")
	fmt.Println("  5. Edit the layout/app.go file to start building your application.")
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

// copyWasmExec finds and copies the wasm_exec.js file.
func copyWasmExec() error {
	goRoot := runtime.GOROOT()
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

// --- Templates ---
func mainGoTemplate(projectName string) string {
	return fmt.Sprintf(`
//go:build js && wasm

package main

import (
	"%s/layout"

	"fmt"
	"github.com/AureClai/vortex/renderer"
)

func main() {
	fmt.Println("Vortex app initialized! 🚀")

	// Create renderer
	r := renderer.NewRenderer("app")

	// Create the welcome page
	app := layout.NewApp(r)

	// Render the app
	r.Render(app.Render())

	// Keep the program running
	<-make(chan bool)
}
	`, projectName)
}

const appGoTemplate = `
//go:build js && wasm

package layout

import (
	"fmt"
	"syscall/js"

	"github.com/AureClai/vortex/component"
	"github.com/AureClai/vortex/renderer"
	"github.com/AureClai/vortex/style"
	"github.com/AureClai/vortex/vdom"
)

// --- Styles ---
var AppStyle = style.New(
	style.CustomStyle("max-width", "800px"),
	style.Margin(style.MarginAll, "0 auto"),
	style.Padding(style.PaddingAll, "2rem"),
	style.CustomStyle("min-height", "100vh"),
	style.Display(style.DisplayFlex),
	style.FlexDirection(style.FlexDirectionColumn),
	style.JustifyContent(style.JustifyContentCenter),
	style.MediaQuery(
		style.MediaQueryTypeMaxWidth,
		"768px",
		style.Padding(style.PaddingAll, "1rem"),
	),
)

var headerStyle = style.New(
	style.TextAlign(style.TextAlignCenter),
	style.Margin(style.MarginAll, "0 auto"),
	style.Padding(style.PaddingAll, "3rem"),
	style.CustomStyle("min-height", "100vh"),
	style.Display(style.DisplayFlex),
	style.FlexDirection(style.FlexDirectionColumn),
	style.JustifyContent(style.JustifyContentCenter),
	style.MediaQuery(
		style.MediaQueryTypeMaxWidth,
		"768px",
		style.Padding(style.PaddingAll, "2rem"),
	),
)

// TODO: Animation will be defined with Vortex Animation Engine
var logoStyle = style.New(
	style.FontSize("4rem"),
	style.Margin(style.MarginBottom, "1rem"),
	style.Display(style.DisplayBlock),
	//style.CustomStyle("animation", "rotate 10s linear infinite"),
	style.MediaQuery(
		style.MediaQueryTypeMaxWidth,
		"768px",
		style.FontSize("3rem"),
	),
)

//TODO:
// @keyframes rotate {
// 	from { transform: rotate(0deg); }
// 	to { transform: rotate(360deg); }
// }

var mainTitleStyle = style.New(
	style.FontSize("3rem"),
	style.FontWeight("700"),
	style.Margin(style.MarginBottom, "1rem"),
	style.Color("#ffffff"),
	style.CustomStyle("text-shadow", "0 4px 20px rgba(0, 0, 0, 0.3)"),
	style.MediaQuery(
		style.MediaQueryTypeMaxWidth,
		"768px",
		style.FontSize("2.2rem"),
	),
)

var subtitleStyle = style.New(
	style.FontSize("1.2rem"),
	style.Color("#ffffff"),
	style.Margin(style.MarginBottom, "0"),
)

// TODO: Add a theme for all the colors used in the app to replace the "var(--shadow)"
var featuresSectionStyle = style.New(
	style.BackgroundColor("#ffffff"),
	style.Padding(style.PaddingAll, "2.5rem"),
	style.BorderRadius("20px"),
	style.Margin(style.MarginBottom, "2rem"),
	style.BoxShadow("0", "10px", "0", "40px", "#000000", false),
)

var sectionTitleStyle = style.New(
	style.Margin(style.MarginBottom, "1.5rem"),
	style.FontWeight("600"),
	style.FontSize("2rem"),
	style.TextAlign(style.TextAlignCenter),
)

var sectionLinksTitleStyle = style.New(
	style.Margin(style.MarginBottom, "1.5rem"),
	style.FontWeight("600"),
	style.FontSize("1.5rem"),
	style.TextAlign(style.TextAlignCenter),
)

var linksSectionStyle = style.New(
	style.TextAlign(style.TextAlignCenter),
)

var buttonStyle = style.New(
	style.Display(style.DisplayInlineBlock),
	style.Padding(style.PaddingAll, "12px 24px"),
	style.Margin(style.MarginAll, "0.5rem"),
	style.Border("none"),
	style.BorderRadius("12px"),
	style.FontWeight("600"),
	style.FontSize("1rem"),
	style.Cursor(style.CursorPointer),
	style.CustomStyle("transition", "all 0.3s ease"),
	style.TextDecoration("none"),
	style.FontFamily("inherit"),
	style.OnHover(
		style.CustomStyle("transform", "translateY(-2px)"),
		style.CustomStyle("box-shadow", "0 8px 25px #000000"),
	),
	style.MediaQuery(
		style.MediaQueryTypeMaxWidth,
		"768px",
		style.Display(style.DisplayBlock),
		style.Margin(style.MarginAll, "0.5rem 0"),
	),
)

var primaryButtonStyle = style.Extend(buttonStyle,
	style.CustomStyle("background", "linear-gradient(45deg, #007bff, #00bfff)"),
	style.Color("#ffffff"),
	style.CustomStyle("box-shadow", "0 4px 15px rgba(33, 150, 243, 0.3)"),
)

var secondaryButtonStyle = style.Extend(buttonStyle,
	style.BackgroundColor("#f0f0f0"),
	style.Color("#000000"),
	style.Border("2px solid #007bff"),
)

var accentButtonStyle = style.Extend(buttonStyle,
	style.BackgroundColor("linear-gradient(45deg, #ff6b6b, #ff8a8a)"),
	style.Color("#ffffff"),
	style.CustomStyle("box-shadow", "0 4px 15px rgba(255, 107, 107, 0.3)"),
)

var footerStyle = style.New(
	style.TextAlign(style.TextAlignCenter),
	style.Padding(style.PaddingAll, "2rem"),
	style.BackgroundColor("rgba(255, 255, 255, 0.1)"),
	style.BorderRadius("20px"),
	style.Border("1px solid rgba(255, 255, 255, 0.2)"),
	style.CustomStyle("backdrop-filter", "blur(20px)"),
)

var footerTextStyle = style.New(
	style.Color("rgba(255, 255, 255, 0.9)"),
	style.FontSize("0.9rem"),
	style.Margin(style.MarginAll, "0"),
)

var ErrorStyle = style.New(
	style.Display(style.DisplayFlex),
	style.FlexDirection(style.FlexDirectionColumn),
	style.JustifyContent(style.JustifyContentCenter),
	style.AlignItems(style.AlignItemsCenter),
	style.Height("100vh"),
	style.Color("white"),
	style.TextAlign(style.TextAlignCenter),
)

type AppState struct {
}

type App struct {
	vdom.StatefulComponentBase[AppState]
}

func NewApp(r *renderer.Renderer) *App {
	app := &App{}

	reRender := func() {
		r.Render(app.Render())
	}

	initialState := AppState{}

	app.StatefulComponentBase = vdom.NewStatefulComponent("div", initialState, reRender)

	return app
}

func (a *App) Render() *vdom.VNode {
	page := component.NewContainer()
	page.Style(AppStyle)

	// Header section
	header := component.NewContainer()
	header.Style(headerStyle)

	// Logo and title
	logo := component.NewText("🔄")
	logo.Style(logoStyle)
	title := component.NewHeading("Welcome to Vortex!", 1)
	title.Style(mainTitleStyle)
	subtitle := component.NewParagraph("Your Go WebAssembly application is running successfully!")
	subtitle.Style(subtitleStyle)
	header.AddChild(logo)
	header.AddChild(title)
	header.AddChild(subtitle)

	// Links section
	links := component.NewContainer()
	links.Style(linksSectionStyle)
	linksTitle := component.NewHeading("Learn More", 3)
	linksTitle.Style(sectionLinksTitleStyle)
	links.AddChild(linksTitle)

	// Create interactive buttons
	docsBtn := component.NewButton("📚 Documentation")
	docsBtn.Style(primaryButtonStyle)
	docsBtn.On("click", func(e js.Value) {
		fmt.Println("Opening Vortex documentation...")
		// In a real app, this would open the docs
	})
	githubBtn := component.NewButton("🐙 GitHub")
	githubBtn.Style(secondaryButtonStyle)
	githubBtn.On("click", func(e js.Value) {
		fmt.Println("Opening Vortex GitHub repository...")
		// In a real app, this would open the GitHub repository
	})
	links.AddChild(docsBtn)
	links.AddChild(githubBtn)

	// Footer
	footer := component.NewContainer()
	footer.Style(footerStyle)
	footerText := component.NewParagraph("Built with ❤️ using Vortex - The Go WebAssembly Framework")
	footerText.Style(footerTextStyle)
	footer.AddChild(footerText)

	// Assemble the page
	page.AddChild(header)
	page.AddChild(links)
	page.AddChild(footer)

	return page.Render()
}
`

const indexHTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vortex App</title>
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
	</style>
	</head>
<body>
    <div id="app">
    </div>

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

func goModTemplate(projectName string) string {
	return fmt.Sprintf(`module %s

go 1.25.1

require github.com/AureClai/vortex v0.1.0
`, projectName)
}

func readmeTemplate(projectName string) string {
	return fmt.Sprintf(`# %s

A Vortex WebAssembly application generated by the Vortex CLI.

## 🚀 Getting Started

### Prerequisites
- Go 1.25.1 or later
- Modern web browser with WebAssembly support

### Development

1. **Build the application:**
`+"   ```bash\n   vortex build\n   ```"+`

2. **Start the development server:**
`+"   ```bash\n   vortex dev\n   ```"+`

3. **Open your browser:**
   Visit http://localhost:8080 to see your app!

4. **Edit the layout/app.go file to start building your application.**

## 📁 Project Structure

- `+"`main.go`"+` - Your main application code
- `+"`index.html`"+` - HTML template with styling
- `+"`go.mod`"+` - Go module configuration
- `+"`app.wasm`"+` - Compiled WebAssembly binary (generated)
- `+"`wasm_exec.js`"+` - Go WebAssembly runtime (generated)
- `+"`layout/app.go`"+` - Your main application code

## 📚 Learn More

- [Vortex Documentation](https://github.com/AureClai/vortex)
- [Vortex Examples](https://github.com/AureClai/vortex/tree/main/examples)
- [WebAssembly with Go](https://golang.org/pkg/syscall/js/)

## 🛠️ Next Steps

1. Edit `+"`main.go`"+` to customize your application
2. Add more Vortex components from the library or create your own
3. Style your app with the Vortex CSS-in-Go API
4. Build something amazing with Go and WebAssembly!

---

Built with ❤️ using [Vortex](https://github.com/AureClai/vortex) - The Go WebAssembly Framework
`, projectName)
}

const gitignoreTemplate = `
# Compiled Wasm file
app.wasm

# JS glue file
wasm_exec.js
`
