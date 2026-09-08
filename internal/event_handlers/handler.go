package event_handlers

import (
	"context"
	"fmt"

	"github.com/google/go-github/v91/github"
	"github.com/tmarback/github-helper-app/internal/commands"
	"golang.org/x/oauth2"
)

// Dynamic provider of token sources for a given installation
type TokenSourceProvider interface {
	TokenSource(ctx context.Context, installation *github.Installation) oauth2.TokenSource
}

// Centralized event handler
type EventHandler struct {
	// Source of application-level authentication tokens
	tokenSourceProvider TokenSourceProvider
	commandRegistry     *commands.CommandRegistry
}

// Creates a new handler that uses the given source of tokens
func NewHandler(tokenSourceProvider TokenSourceProvider) (*EventHandler, error) {

	if tokenSourceProvider == nil {
		return nil, fmt.Errorf("token source provider must be specified")
	}

	return &EventHandler{
		tokenSourceProvider: tokenSourceProvider,
		commandRegistry:     commands.NewRegistry(),
	}, nil

}

// Obtains a Github client for the given installation
func (handler *EventHandler) getGithubClient(ctx context.Context, installation *github.Installation) (*github.Client, error) {

	tokenSource := handler.tokenSourceProvider.TokenSource(ctx, installation)
	httpClient := oauth2.NewClient(ctx, tokenSource)
	client, err := github.NewClient(github.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to obtain github client: %w", err)
	}
	return client, nil

}
