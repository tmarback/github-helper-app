package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v91/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

// Github App-based token source provider
type GithubAppTokenSourceProvider struct {
	// Underlying application token source
	ApplicationTokenSource oauth2.TokenSource
}

func (provider *GithubAppTokenSourceProvider) TokenSource(ctx context.Context, installation *github.Installation) oauth2.TokenSource {
	return githubauth.NewInstallationTokenSource(*installation.ID, provider.ApplicationTokenSource)
}

// Configuration for Github App provider
type GithubAppTokenConfig struct {
	// Path to file that contains the private key
	PrivateKeyPath string `koanf:"privateKeyPath"`
	// App client ID
	ClientId string `koanf:"clientId"`
}

// Creates a new instance
func (config *GithubAppTokenConfig) Create() (*GithubAppTokenSourceProvider, error) {

	// Read private key
	if config.PrivateKeyPath == "" {
		return nil, fmt.Errorf("private key path must be provided")
	}
	privateKey, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Create application token source
	if config.ClientId == "" {
		return nil, fmt.Errorf("client ID must be provided")
	}
	applicationTokenSource, err := githubauth.NewApplicationTokenSource(config.ClientId, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating application token source: %w", err)
	}

	return &GithubAppTokenSourceProvider{
		ApplicationTokenSource: applicationTokenSource,
	}, nil

}
