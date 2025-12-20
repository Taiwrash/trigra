// Package github provides the GitHub implementation of the Provider interface.
package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/google/go-github/v79/github"
)

// Provider implements the providers.Provider interface for GitHub.
type Provider struct {
	client *github.Client
}

// NewProvider creates a new GitHub provider instance.
func NewProvider(token string) *Provider {
	var client *github.Client
	if token != "" {
		client = github.NewClient(nil).WithAuthToken(token)
	} else {
		client = github.NewClient(nil)
	}
	return &Provider{client: client}
}

// Name returns "github".
func (p *Provider) Name() string {
	return "github"
}

// Validate validates the GitHub webhook payload.
func (p *Provider) Validate(r *http.Request, secret string) ([]byte, error) {
	return github.ValidatePayload(r, []byte(secret))
}

// ParsePushEvent parses a GitHub push event payload.
func (p *Provider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, err
	}

	pushEvent, ok := event.(*github.PushEvent)
	if !ok {
		return nil, fmt.Errorf("event is not a push event")
	}

	modifiedFiles := make(map[string]bool)
	// Suppress SA1019: Commits is deprecated but used for efficiency in the webhook payload.
	//nolint:staticcheck
	for _, commit := range pushEvent.Commits {
		for _, file := range commit.Added {
			modifiedFiles[file] = true
		}
		for _, file := range commit.Modified {
			modifiedFiles[file] = true
		}
	}

	files := make([]string, 0, len(modifiedFiles))
	for f := range modifiedFiles {
		files = append(files, f)
	}

	return &providers.PushEvent{
		Owner:         pushEvent.Repo.GetOwner().GetName(),
		Repo:          pushEvent.Repo.GetName(),
		Ref:           pushEvent.GetRef(),
		After:         pushEvent.GetAfter(),
		ModifiedFiles: files,
	}, nil
}

// DownloadFile downloads a file from GitHub.
func (p *Provider) DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
	fileReader, _, err := p.client.Repositories.DownloadContents(
		ctx,
		owner,
		repo,
		path,
		&github.RepositoryContentGetOptions{
			Ref: ref,
		},
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fileReader.Close() }()

	return io.ReadAll(fileReader)
}
