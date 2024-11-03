package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/nathan-fiscaletti/kbm-overlay/internal/config"
	"github.com/nathan-fiscaletti/kbm-overlay/internal/listener"
	"github.com/nathan-fiscaletti/kbm-overlay/internal/relay"
)

var outChan = make(chan any)

func main() {
	// Make sure we are running as root
	if os.Geteuid() != 0 {
		fmt.Println("Error: This program must be run as root.")
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	exeDir := filepath.Dir(exePath)

	cfg, err := config.LoadConfig(filepath.Join(exeDir, "config.yaml"))
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	fmt.Println("Loaded config: config.yaml")
	fmt.Println()

	// Configure wsurl.js
	webDir := filepath.Join(exeDir, "web")
	wsUrl := fmt.Sprintf("ws://localhost:%d/ws", cfg.Port)
	err = os.WriteFile(filepath.Join(webDir, "wsurl.js"), []byte(fmt.Sprintf("const wsUrl = \"%s\";\n", wsUrl)), 0644)
	if err != nil {
		fmt.Println("Error writing wsurl.js:", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = listener.Listen(ctx, cfg, outChan)
	if err != nil {
		fmt.Println("Error listening for events:", err)
		return
	}

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		// Get the file name from the URL path
		file := path.Base(r.URL.Path)

		// If the file name is "view" or empty (when requesting /view/), serve "index.html" by default
		if file == "view" || file == "" {
			file = "index.html"
		}

		// Define the path to the static directory
		filePath := path.Join(webDir, file)

		// Serve the file
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/ws", relay.NewWsHandler(outChan))

	viewUrl := fmt.Sprintf("http://localhost:%d/view", cfg.Port)

	fmt.Printf("Overlay available at %s\n", viewUrl)

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
