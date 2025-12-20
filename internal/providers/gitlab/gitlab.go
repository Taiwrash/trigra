// Package gitlab provides the GitLab implementation of the Provider interface.
package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/xanzy/go-gitlab"
)

// Provider implements the providers.Provider interface for GitLab.
type Provider struct {
	//nolint:staticcheck
	client *gitlab.Client
}

// NewProvider creates a new GitLab provider instance.
func NewProvider(token string) *Provider {
	//nolint:staticcheck
	client, _ := gitlab.NewClient(token)
	return &Provider{client: client}
}

// Name returns "gitlab".
func (p *Provider) Name() string {
	return "gitlab"
}

// Validate validates the GitLab webhook payload.
func (p *Provider) Validate(r *http.Request, secret string) ([]byte, error) {
	// GitLab uses X-Gitlab-Token header for secret validation
	receivedSecret := r.Header.Get("X-Gitlab-Token")
	if receivedSecret != "" && receivedSecret != secret {
		return nil, fmt.Errorf("invalid gitlab secret")
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// ParsePushEvent parses a GitLab push event payload.
func (p *Provider) ParsePushEvent(_ *http.Request, payload []byte) (*providers.PushEvent, error) {
	var event gitlab.PushEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}

	if event.ObjectKind != "push" {
		return nil, fmt.Errorf("not a push event")
	}

	modifiedFiles := make(map[string]bool)
	for _, commit := range event.Commits {
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
		Owner:         event.Project.Namespace,
		Repo:          event.Project.Name,
		Ref:           event.Ref,
		After:         event.After,
		ModifiedFiles: files,
	}, nil
}

// DownloadFile downloads a file from GitLab.
func (p *Provider) DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
	// Project ID can be the path with namespace
	projectID := fmt.Sprintf("%s/%s", owner, repo)
	file, _, err := p.client.RepositoryFiles.GetRawFile(projectID, path, &gitlab.GetRawFileOptions{
		Ref: gitlab.Ptr(ref),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// SetupWebhook ensures a GitLab webhook is configured.
func (p *Provider) SetupWebhook(ctx context.Context, owner, repo, url, secret string) error {
	projectID := fmt.Sprintf("%s/%s", owner, repo)

	hooks, _, err := p.client.Projects.ListProjectHooks(projectID, &gitlab.ListProjectHooksOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to list project hooks: %w", err)
	}

	for _, hook := range hooks {
		if hook.URL == url {
			return nil // Hook already exists
		}
	}

	opt := &gitlab.AddProjectHookOptions{
		URL:        gitlab.Ptr(url),
		PushEvents: gitlab.Ptr(true),
		Token:      gitlab.Ptr(secret),
	}

	_, _, err = p.client.Projects.AddProjectHook(projectID, opt, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to add project hook: %w", err)
	}

	return nil
}
