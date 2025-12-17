package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the GitOps controller
type Config struct {
	// GitHub configuration
	GitHubToken   string
	WebhookSecret string

	// Server configuration
	ServerPort int

	// Kubernetes configuration
	InCluster bool
	Namespace string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		GitHubToken:   os.Getenv("GITHUB_TOKEN"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		ServerPort:    8082, // default
		Namespace:     getEnvOrDefault("NAMESPACE", "default"),
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

	// GitHub token is optional for public repos, but recommended
	if cfg.GitHubToken == "" {
		fmt.Println("WARNING: GITHUB_TOKEN not set. Rate limits will be lower for GitHub API.")
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
