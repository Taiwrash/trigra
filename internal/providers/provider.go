// Package providers defines the Git provider abstractions for TRIGRA.
package providers

import (
	"context"
	"net/http"
)

// PushEvent represents a generic push event from any provider
type PushEvent struct {
	Owner         string
	Repo          string
	Ref           string
	After         string // commit SHA
	ModifiedFiles []string
}

// Provider defines the interface that all Git providers must implement
type Provider interface {
	// Name returns the name of the provider (e.g., "github", "gitlab")
	Name() string

	// Validate validates the webhook payload using the provided secret
	Validate(r *http.Request, secret string) ([]byte, error)

	// ParsePushEvent parses the provider-specific payload into a generic PushEvent
	ParsePushEvent(r *http.Request, payload []byte) (*PushEvent, error)

	// DownloadFile downloads a single file from the repository
	DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error)
}
