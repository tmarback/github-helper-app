package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// Server configuration
type Config struct {
	// The minimum level to log
	LogLevel string `koanf:"logLevel"`
	// The secret that the webhook must have to be accepted
	WebhookSecret string `koanf:"webhookSecret"`
	// The hostname to bind to
	Hostname string `koanf:"hostname"`
	// The port to bind to
	Port int `koanf:"port"`
	// Configuration of auth via Github App
	GithubAppAuth *GithubAppTokenConfig `koanf:"githubAuth"`
}

// Loads the server configuration
func LoadConfig() (config Config, err error) {

	var k = koanf.New(".")

	// Load defaults
	if err := k.Load(structs.Provider(Config{
		WebhookSecret: "",
		Hostname:      "0.0.0.0",
		Port:          8000,
	}, "koanf"), nil); err != nil {
		return config, fmt.Errorf("failed to load defaults: %w", err)
	}

	// Load from config file if present
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		// If the file simply doesn't exist, we skip it and use defaults
		if errors.Is(err, os.ErrNotExist) {
			slog.Debug("No config file found")
		} else {
			return config, fmt.Errorf("failed to load config file: %w", err)
		}

	}

	// Load from environment
	const ENV_PREFIX = "CONFIG__"
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: ENV_PREFIX,
		TransformFunc: func(k, v string) (string, any) {
			slog.Debug("Processing key " + k)

			// Remove prefix
			k = strings.TrimPrefix(k, ENV_PREFIX)

			// Obtain the key segments
			segments := strings.Split(k, "__")

			// Normalize each segment to camel-case
			for i, seg := range segments {
				segments[i] = strcase.ToLowerCamel(seg)
			}

			// Create key
			k = strings.Join(segments, ".")
			slog.Debug("Processed key " + k)

			return k, v
		},
	}), nil); err != nil {
		return config, fmt.Errorf("failed to load environment: %w", err)
	}

	// Unmarshal config into struct
	if err := k.Unmarshal("", &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return

}
