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

type GitLabProvider struct {
	client *gitlab.Client
}

func NewGitLabProvider(token string) *GitLabProvider {
	client, _ := gitlab.NewClient(token)
	return &GitLabProvider{client: client}
}

func (p *GitLabProvider) Name() string {
	return "gitlab"
}

func (p *GitLabProvider) Validate(r *http.Request, secret string) ([]byte, error) {
	// GitLab uses X-Gitlab-Token header for secret validation
	receivedSecret := r.Header.Get("X-Gitlab-Token")
	if receivedSecret != secret {
		return nil, fmt.Errorf("invalid gitlab secret")
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (p *GitLabProvider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
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

func (p *GitLabProvider) DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
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
