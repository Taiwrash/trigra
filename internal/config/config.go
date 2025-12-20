package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the GitOps controller
type Config struct {
	// Git provider configuration
	GitProvider   string // "github", "gitlab", etc.
	GitToken      string // Token for the git provider API
	WebhookSecret string // Secret to validate webhooks

	// Server configuration
	ServerPort int

	// Kubernetes configuration
	InCluster bool
	Namespace string

	// Generic Git configuration
	GitRepoURL string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		GitProvider:   getEnvOrDefault("GIT_PROVIDER", "github"),
		GitToken:      getEnvOrDefault("GIT_TOKEN", os.Getenv("GITHUB_TOKEN")),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		ServerPort:    8082, // default
		Namespace:     getEnvOrDefault("NAMESPACE", "default"),
		GitRepoURL:    os.Getenv("GIT_REPO_URL"),
	}

	// Parse server port if provided
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
		}
		cfg.ServerPort = port
	}

	// Auto-detect in-cluster mode
	cfg.InCluster = isInCluster()

	// Validate required fields
	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required configuration is present
func Validate(cfg *Config) error {
	if cfg.WebhookSecret == "" {
		return fmt.Errorf("WEBHOOK_SECRET is required")
	}

	// Git token is optional for some public repos, but recommended
	if cfg.GitToken == "" {
		fmt.Printf("WARNING: GIT_TOKEN not set for provider %s. API access might be restricted.\n", cfg.GitProvider)
	}

	return nil
}

// isInCluster detects if we're running inside a Kubernetes cluster
func isInCluster() bool {
	// Check for service account token (standard in-cluster indicator)
	_, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token")
	return err == nil
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
