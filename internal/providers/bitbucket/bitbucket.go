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

// BitbucketProvider implements the providers.Provider interface for Bitbucket.
type BitbucketProvider struct {
	client *bitbucket.Client
}

// NewBitbucketProvider creates a new Bitbucket provider instance.
func NewBitbucketProvider(user, token string) *BitbucketProvider {
	client, _ := bitbucket.NewBasicAuth(user, token)
	return &BitbucketProvider{client: client}
}

// Name returns "bitbucket".
func (p *BitbucketProvider) Name() string {
	return "bitbucket"
}

// Validate validates the Bitbucket webhook payload.
func (p *BitbucketProvider) Validate(r *http.Request, _ string) ([]byte, error) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// ParsePushEvent parses a Bitbucket push event payload.
func (p *BitbucketProvider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
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
func (p *BitbucketProvider) DownloadFile(_ context.Context, owner, repo, ref, path string) ([]byte, error) {
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
