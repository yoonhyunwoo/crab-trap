package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/your-username/moltbook-prompt-injector/internal/config"
	"github.com/your-username/moltbook-prompt-injector/internal/server"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to config file")
	port       = flag.Int("port", 8080, "Server port")
	logDir     = flag.String("log-dir", "./logs", "Log directory")
)

func main() {
	flag.Parse()

	cfg, err := config.LoadDefault()
	if err != nil {
		log.Printf("Warning: Failed to load config, using defaults: %v", err)
		cfg = &config.Config{
			Server: config.ServerConfig{
				Port:   *port,
				LogDir: *logDir,
			},
		}
	}

	logger, err := server.NewLogger(cfg.Server.LogDir)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	handler := server.NewHandler(logger)

	http.HandleFunc("/", handler.HandleUI)
	http.HandleFunc("/log", handler.HandleRequest)
	http.HandleFunc("/health", handler.HandleHealth)
	http.HandleFunc("/logs", handler.HandleLogs)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on http://localhost%s", addr)
	log.Printf("Logs directory: %s", cfg.Server.LogDir)
	log.Printf("Web UI: http://localhost%s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
