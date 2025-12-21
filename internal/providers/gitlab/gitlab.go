// Package gitlab provides the GitLab implementation of the Provider interface.
package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Taiwrash/trigra/internal/providers"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Provider implements the providers.Provider interface for GitLab.
type Provider struct {
	client *gitlab.Client
}

// NewProvider creates a new GitLab provider instance.
func NewProvider(baseURL, token string) *Provider {
	var client *gitlab.Client
	var err error

	if baseURL != "" {
		client, err = gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	} else {
		client, err = gitlab.NewClient(token)
	}

	if err != nil {
		log.Printf("ERROR: Failed to create GitLab client: %v", err)
	}

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
	if secret != "" && receivedSecret != secret {
		return nil, fmt.Errorf("invalid or missing gitlab secret")
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// ParsePushEvent parses a GitLab push event payload.
func (p *Provider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
	eventKey := r.Header.Get("X-Gitlab-Event")
	if eventKey != "" && eventKey != "Push Hook" && eventKey != "Tag Push Hook" {
		return nil, fmt.Errorf("ignoring gitlab event: %s", eventKey)
	}

	var event gitlab.PushEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}

	// Some payloads don't have ObjectKind set but are still push events
	if event.ObjectKind != "" && event.ObjectKind != "push" && event.ObjectKind != "tag_push" {
		return nil, fmt.Errorf("not a push event: %s", event.ObjectKind)
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

	// Use PathWithNamespace to handle groups/subgroups correctly
	parts := strings.Split(event.Project.PathWithNamespace, "/")
	owner := strings.Join(parts[:len(parts)-1], "/")
	repo := parts[len(parts)-1]

	return &providers.PushEvent{
		Owner:         owner,
		Repo:          repo,
		Ref:           event.Ref,
		After:         event.After,
		ModifiedFiles: files,
	}, nil
}

// DownloadFile downloads a file from GitLab.
func (p *Provider) DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
	// Project ID is the path with namespace
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
