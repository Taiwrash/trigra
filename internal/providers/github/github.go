package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/google/go-github/v79/github"
)

type GitHubProvider struct {
	client *github.Client
}

func NewGitHubProvider(token string) *GitHubProvider {
	var client *github.Client
	if token != "" {
		client = github.NewClient(nil).WithAuthToken(token)
	} else {
		client = github.NewClient(nil)
	}
	return &GitHubProvider{client: client}
}

func (p *GitHubProvider) Name() string {
	return "github"
}

func (p *GitHubProvider) Validate(r *http.Request, secret string) ([]byte, error) {
	return github.ValidatePayload(r, []byte(secret))
}

func (p *GitHubProvider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, err
	}

	pushEvent, ok := event.(*github.PushEvent)
	if !ok {
		return nil, fmt.Errorf("event is not a push event")
	}

	modifiedFiles := make(map[string]bool)
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

func (p *GitHubProvider) DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
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
	defer fileReader.Close()

	return io.ReadAll(fileReader)
}
