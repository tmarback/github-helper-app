package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/cbrgm/githubevents/v2/githubevents"
	"github.com/tmarback/github-helper-app/internal/event_handlers"
)

// Parse the requested log level
func parseLogLevel(s string) (level slog.Level, err error) {

	err = level.UnmarshalText([]byte(s))
	return

}

func main() {

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load application config: %v", err)
	}

	// Set up logging
	logLevel, err := parseLogLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}
	log.Printf("Using log level %v", logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	var tokenSourceProvider event_handlers.TokenSourceProvider
	if config.GithubAppAuth != nil {
		slog.Info("Configuring Github App authentication")
		tokenSourceProvider, err = config.GithubAppAuth.Create()
		if err != nil {
			log.Fatalf("Failed to create app auth provider: %v", err)
		}
	} else {
		log.Fatalf("No auth configuration was provided")
	}

	// Create event handler
	eventHandler, err := event_handlers.NewHandler(tokenSourceProvider)
	if err != nil {
		log.Fatalf("Failed to create event handler: %v", err)
	}

	// Initialize event manager
	if config.WebhookSecret == "" {
		slog.Warn("Webhook secret not configured")
	}
	handle := githubevents.New(config.WebhookSecret)

	// Register handlers
	handle.OnIssueCommentCreated(eventHandler.HandleIssueComment)

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
