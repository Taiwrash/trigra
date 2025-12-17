package github

import (
	"context"

	"github.com/google/go-github/v79/github"
	"golang.org/x/oauth2"
)

// NewClient creates a new GitHub client with optional authentication
func NewClient(ctx context.Context, token string) *github.Client {
	if token == "" {
		// Return unauthenticated client (lower rate limits)
		return github.NewClient(nil)
	}

	// Create authenticated client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
