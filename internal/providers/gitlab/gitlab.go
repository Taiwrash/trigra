package gitlab

import (
	"context"
	"errors"
	"net/http"

	"github.com/Taiwrash/trigra/internal/providers"
)

type GitLabProvider struct {
	token string
}

func NewGitLabProvider(token string) *GitLabProvider {
	return &GitLabProvider{token: token}
}

func (p *GitLabProvider) Name() string {
	return "gitlab"
}

func (p *GitLabProvider) Validate(r *http.Request, secret string) ([]byte, error) {
	// TODO: Implement GitLab webhook secret validation (X-Gitlab-Token)
	return nil, errors.New("gitlab provider not yet fully implemented")
}

func (p *GitLabProvider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
	// TODO: Implement GitLab push event parsing
	return nil, errors.New("gitlab provider not yet fully implemented")
}

func (p *GitLabProvider) DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
	// TODO: Implement GitLab API file download
	return nil, errors.New("gitlab provider not yet fully implemented")
}
