package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/cbrgm/githubevents/v2/githubevents"
	"github.com/tmarback/github-helper-app/internal/event_handlers"
)

func main() {

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load application config: %v", err)
	}

	// Initialize event manager
	if config.WebhookSecret == "" {
		slog.Warn("Webhook secret not configured")
	}
	handle := githubevents.New(config.WebhookSecret)

	// Register handlers
	handle.OnIssueCommentCreated(event_handlers.HandleIssueComment)

	// Configure HTTP endpoint
	http.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		err := handle.HandleEventRequest(r)
		if err != nil {
			slog.Error("Error during event handling", slog.Any("error", err))
		}
	})

	// Start HTTP server
	addr := fmt.Sprintf("%s:%d", config.Hostname, config.Port)
	slog.Info("Starting server", slog.String("address", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panic(err)
	}

}
