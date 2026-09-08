package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/cbrgm/githubevents/v2/githubevents"
	"github.com/jferrl/go-githubauth"
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

	// Read private key
	if config.GithubAuth.PrivateKeyPath == "" {
		log.Fatalf("Private key path must be provided")
	}
	privateKey, err := os.ReadFile(config.GithubAuth.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}

	// Create application token source
	applicationTokenSource, err := githubauth.NewApplicationTokenSource(config.GithubAuth.ClientId, privateKey)
	if err != nil {
		log.Fatalf("Error creating application token source: %v", err)
	}

	// Create event handler
	eventHandler := event_handlers.NewHandler(applicationTokenSource)

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
