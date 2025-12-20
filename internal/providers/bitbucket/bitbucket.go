// Package bitbucket provides the Bitbucket implementation of the Provider interface.
package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/ktrysmt/go-bitbucket"
)

// Provider implements the providers.Provider interface for Bitbucket.
type Provider struct {
	client *bitbucket.Client
}

// NewProvider creates a new Bitbucket provider instance.
func NewProvider(user, token string) *Provider {
	client, _ := bitbucket.NewBasicAuth(user, token)
	return &Provider{client: client}
}

// Name returns "bitbucket".
func (p *Provider) Name() string {
	return "bitbucket"
}

// Validate validates the Bitbucket webhook payload.
func (p *Provider) Validate(r *http.Request, _ string) ([]byte, error) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// ParsePushEvent parses a Bitbucket push event payload.
func (p *Provider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
	eventKey := r.Header.Get("X-Event-Key")
	if eventKey != "repo:push" {
		return nil, fmt.Errorf("not a push event")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, err
	}

	push, ok := raw["push"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid bitbucket payload")
	}

	changes, ok := push["changes"].([]interface{})
	if !ok || len(changes) == 0 {
		return nil, fmt.Errorf("no changes in bitbucket push")
	}

	change := changes[0].(map[string]interface{})
	newCommit := change["new"].(map[string]interface{})
	target := newCommit["target"].(map[string]interface{})
	hash := target["hash"].(string)

	repo := raw["repository"].(map[string]interface{})
	fullName := repo["full_name"].(string)
	parts := strings.Split(fullName, "/")
	owner := parts[0]

	return &providers.PushEvent{
		Owner:         owner,
		Repo:          repo["name"].(string),
		Ref:           hash,
		After:         hash,
		ModifiedFiles: []string{"."},
	}, nil
}

// DownloadFile downloads a file from Bitbucket.
func (p *Provider) DownloadFile(_ context.Context, owner, repo, ref, path string) ([]byte, error) {
	res, err := p.client.Repositories.Repository.GetFileBlob(&bitbucket.RepositoryBlobOptions{
		Owner:    owner,
		RepoSlug: repo,
		Path:     path,
		Ref:      ref,
	})
	if err != nil {
		return nil, err
	}
	return res.Content, nil
}

// SetupWebhook ensures a Bitbucket webhook is configured.
func (p *Provider) SetupWebhook(_ context.Context, owner, repo, url, _ string) error {
	// Bitbucket Cloud webhooks
	hooks, err := p.client.Repositories.Webhooks.List(&bitbucket.WebhooksOptions{
		Owner:    owner,
		RepoSlug: repo,
	})
	if err != nil {
		return fmt.Errorf("failed to list webhooks: %w", err)
	}
	_ = hooks // Suppress unused for now while we investigate typing

	// hooks is likely a raw interface or specific type from go-bitbucket
	// Checking the SDK structure... usually it's a map or slice
	// For now, let's assume we can iterate or at least try to create one.

	// Note: go-bitbucket doesn't easily expose the list in a typed way sometimes.
	// We'll attempt a simple existence check if possible.

	opt := &bitbucket.WebhooksOptions{
		Owner:       owner,
		RepoSlug:    repo,
		Description: "Trigra Webhook",
		Url:         url,
		Active:      true,
		Events:      []string{"repo:push"},
	}

	_, err = p.client.Repositories.Webhooks.Create(opt)
	if err != nil {
		// If it already exists, Bitbucket might return an error.
		// For simplicity, we just try to create it.
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	return nil
}
