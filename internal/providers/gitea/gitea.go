// Package gitea provides the Gitea implementation of the Provider interface.
package gitea

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.gitea.io/sdk/gitea"
	"github.com/Taiwrash/trigra/internal/providers"
)

// Provider implements the providers.Provider interface for Gitea.
type Provider struct {
	client *gitea.Client
}

// PushPayload defines the structure of a Gitea push event payload.
type PushPayload struct {
	Ref     string `json:"ref"`
	Before  string `json:"before"`
	After   string `json:"after"`
	Commits []struct {
		ID       string   `json:"id"`
		Message  string   `json:"message"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
	Repo struct {
		Name  string `json:"name"`
		Owner struct {
			UserName string `json:"username"`
		} `json:"owner"`
	} `json:"repository"`
}

// NewProvider creates a new Gitea provider instance.
func NewProvider(baseURL, token string) *Provider {
	client, _ := gitea.NewClient(baseURL, gitea.SetToken(token))
	return &Provider{client: client}
}

// Name returns "gitea".
func (p *Provider) Name() string {
	return "gitea"
}

// Validate validates the Gitea webhook payload.
func (p *Provider) Validate(r *http.Request, secret string) ([]byte, error) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if secret != "" {
		signature := r.Header.Get("X-Gitea-Signature")
		if signature == "" {
			return nil, fmt.Errorf("missing Gitea signature")
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(payload)
		expectedSignature := hex.EncodeToString(mac.Sum(nil))

		if signature != expectedSignature {
			return nil, fmt.Errorf("invalid Gitea signature")
		}
	}

	return payload, nil
}

// ParsePushEvent parses a Gitea push event payload.
func (p *Provider) ParsePushEvent(_ *http.Request, payload []byte) (*providers.PushEvent, error) {
	var event PushPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
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
		Owner:         event.Repo.Owner.UserName,
		Repo:          event.Repo.Name,
		Ref:           event.Ref,
		After:         event.After,
		ModifiedFiles: files,
	}, nil
}

// DownloadFile downloads a file from Gitea.
func (p *Provider) DownloadFile(_ context.Context, owner, repo, ref, path string) ([]byte, error) {
	data, _, err := p.client.GetFile(owner, repo, ref, path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SetupWebhook ensures a Gitea webhook is configured.
func (p *Provider) SetupWebhook(_ context.Context, owner, repo, url, secret string) error {
	hooks, _, err := p.client.ListRepoHooks(owner, repo, gitea.ListHooksOptions{})
	if err != nil {
		return fmt.Errorf("failed to list repo hooks: %w", err)
	}

	for _, hook := range hooks {
		if hook.Config["url"] == url {
			return nil // Hook already exists
		}
	}

	_, _, err = p.client.CreateRepoHook(owner, repo, gitea.CreateHookOption{
		Type: gitea.HookTypeGitea,
		Config: map[string]string{
			"url":          url,
			"content_type": "json",
			"secret":       secret,
		},
		Events: []string{"push"},
		Active: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create repo hook: %w", err)
	}

	return nil
}
