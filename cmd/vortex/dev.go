package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// Clients are the connected clients to the websocket
	clients = make(map[*websocket.Conn]bool)
	// Mutex to protect the clients map
	clientsMutex = sync.Mutex{}
	// Rebuild is the channel to rebuild the application
	rebuild = make(chan bool, 1)
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs a local development server for the Vortex application.",
	Long: `Starts a static file server on port 8080 to serve the application.
It is recommended to run 'vortex build' before starting the dev server.`,
	Run: runDev,
}

// printIdleState prints the clean idle state
func printIdleState() {
	// Clear screen (works on Windows, macOS, Linux)
	fmt.Print("\033[2J\033[H")

	fmt.Println("🔨 Vortex development server started")
	fmt.Println(" ")
	fmt.Println("\t🌪️ Vortex app served at: http://localhost:8080")
	fmt.Println("\tWaiting for 🖊️ edition of files to hot reload 🔃")
	fmt.Println(" ")
}

// runDev handles the logic for the 'vortex dev' command.
func runDev(cmd *cobra.Command, args []string) {
	// Initial build
	fmt.Println("🔨 Initial build...")
	if err := buildProject(); err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
		return
	}

	// Start file watcher
	go startFileWatcher()

	// Start rebuild worker
	go rebuildWorker()

	// Start HTTP handlers
	setupHandlers()

	// Show idle state
	printIdleState()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
		return
	}
}

// setupHandlers sets up the HTTP handlers
func setupHandlers() {
	// Websocket endpoint for live reload
	http.HandleFunc("/ws", handleWebSocket)

	// Serve static files with live reload script injection
	http.HandleFunc("/", serveWithReload)
}

// handleWebSocket handles the websocket connection
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade to websocket: %v", err)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Remove client when connection closed
	defer func() {
		fmt.Println("🔄 Client disconnected")
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Don't log "going away" errors since they are expected
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Println("🔄 Client disconnected (expected)")
			} else {
				fmt.Printf("❌ Websocket read error: %v\n", err)
			}
			break
		}
	}
}

// serveWithReload serves files and injects live reload script into HTML
func serveWithReload(w http.ResponseWriter, r *http.Request) {
	// Handle root path
	if r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}

	// Get file path
	filePath := "." + r.URL.Path

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// For HTML files, inject live reload script
	if strings.HasSuffix(filePath, ".html") {
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Inject live reload script beafore closing </body> tag
		htmlContent := string(content)
		liveReloadScript := `
		<script>
(function() {
    console.log('🔄 Vortex Live Reload: Initializing...');
    
    let ws;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 10;
    
    function connect() {
        ws = new WebSocket('ws://localhost:8080/ws');
        
        ws.onopen = function() {
            console.log('🔗 Vortex Live Reload: Connected');
            reconnectAttempts = 0;
        };
        
        ws.onmessage = function(event) {
            if (event.data === 'reload') {
                console.log('🔄 Vortex Live Reload: File changed, reloading...');
                window.location.reload();
            }
        };
        
        ws.onclose = function() {
            console.log('❌ Vortex Live Reload: Connection lost');
            
            // Attempt to reconnect
            if (reconnectAttempts < maxReconnectAttempts) {
                reconnectAttempts++;
                console.log('🔄 Vortex Live Reload: Reconnecting... (attempt ' + reconnectAttempts + ')');
                setTimeout(connect, 1000 * reconnectAttempts);
            } else {
                console.log('🚨 Vortex Live Reload: Max reconnection attempts reached');
            }
        };
        
        ws.onerror = function(error) {
            console.error('🚨 Vortex Live Reload: WebSocket error:', error);
        };
    }
    
    // Initial connection
    connect();
})();
</script>`

		// Insert before closing body tag, or at the end if no body tag found
		if strings.Contains(htmlContent, "</body") {
			htmlContent = strings.Replace(htmlContent, "</body>", liveReloadScript+"</body>", 1)
		} else {
			htmlContent += liveReloadScript
		}

		// Write content to response
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
		return
	}

	http.FileServer(http.Dir(".")).ServeHTTP(w, r)

}

// startFileWatcher watches for changes in Go files
func startFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add directories to watch
	err = addWatchPaths(watcher, ".")
	if err != nil {
		log.Printf("Error setting up file watcher: %v", err)
		return
	}

	fmt.Println("👀 File watcher started")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Debug: Print all the events:
			fmt.Printf("🔍 Event: %s %s\n", event.Op, event.Name)

			// Reacr to write/create operations on Go files
			if (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) && strings.HasSuffix(event.Name, ".go") {
				fmt.Printf("📝 File changed: %s\n", event.Name)

				// Debounce rapid changes
				select {
				case rebuild <- true:
					fmt.Println("🔄 Rebuild triggered")
				default:
					fmt.Println("🔄 Rebuild skipped (debounce)")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)
		}
	}
}

func addWatchPaths(watcher *fsnotify.Watcher, root string) error {
	fmt.Printf("📁 Adding watch paths starting from: %s\n", root)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("❌ Error walking path %s: %v\n", path, err)
			return err
		}

		// Skip hidden directories and common non-source directories
		if info.IsDir() {
			name := filepath.Base(path)
			if name != "." && (strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor") {
				fmt.Printf("⏭️ Skipping directory: %s\n", path)
				return filepath.SkipDir
			}

			fmt.Printf("🔍 Watching directory: %s\n", path)
			if err := watcher.Add(path); err != nil {
				fmt.Printf("❌ Failed to watch %s: %v\n", path, err)
				return err
			}
		}
		return nil
	})
}

func rebuildWorker() {
	for range rebuild {
		// Small delay to prevent rapid rebuilds
		time.Sleep(100 * time.Millisecond) // initial debounce

		// Drain any additional rebuild requests that came in during the debounce
	drainLoop:
		for {
			select {
			case <-rebuild:
				continue // continue draining
			default:
				break drainLoop // Exit the drain loop
			}
		}

		fmt.Println("🔨 Rebuilding...")

		if err := buildProject(); err != nil {
			fmt.Printf("❌ Build failed: %v\n", err)
			continue
		}

		fmt.Println("✅ Build successful, reloading browser...")

		// Notify all connected clients to reload
		clientsMutex.Lock()
		clientCount := len(clients)
		fmt.Printf("📡 Notifying %d connected clients to reload\n", clientCount)

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
			if err != nil {
				fmt.Printf("❌ Failed to send reload message: %v\n", err)
				client.Close()
				delete(clients, client)
			} else {
				fmt.Println("✅ Reload message sent successfully")
			}
		}
		clientsMutex.Unlock()

		fmt.Println("🔄 Hot Reload completed")

		// Wait a moment to let user see the success message
		time.Sleep(1500 * time.Millisecond)

		// Print idle state again
		printIdleState()
	}
}

// buildProject builds the project
func buildProject() error {
	// Check if we are in a vortex project
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		return fmt.Errorf("no main.go file found. Are you in a Vortex project directory?")
	}

	// Build the project
	buildCmd := exec.Command("go", "build", "-o", "app.wasm", ".")
	buildCmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")

	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %s: %w", output, err)
	}

	return nil
}
