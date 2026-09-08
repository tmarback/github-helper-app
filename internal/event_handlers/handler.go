package event_handlers

import (
	"context"
	"fmt"

	"github.com/google/go-github/v91/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

// Centralized event handler
type EventHandler struct {
	// Source of application-level authentication tokens
	applicationTokenSource oauth2.TokenSource
}

// Creates a new handler that uses the given source of application tokens
func NewHandler(applicationTokenSource oauth2.TokenSource) *EventHandler {
	return &EventHandler{
		applicationTokenSource: applicationTokenSource,
	}
}

// Obtains a Github client for the given installation
func (handler *EventHandler) getGithubClient(ctx context.Context, installation *github.Installation) (*github.Client, error) {

	installationTokenSource := githubauth.NewInstallationTokenSource(*installation.ID, handler.applicationTokenSource)
	httpClient := oauth2.NewClient(ctx, installationTokenSource)
	client, err := github.NewClient(github.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to obtain github client: %w", err)
	}
	return client, nil

}
