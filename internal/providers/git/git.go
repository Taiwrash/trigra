package git

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/go-git/go-git/v5"
)

type GenericGitProvider struct {
	baseDir string
	repoURL string
}

func NewGenericGitProvider(repoURL string) *GenericGitProvider {
	tmpDir, _ := os.MkdirTemp("", "trigra-git-*")
	return &GenericGitProvider{
		baseDir: tmpDir,
		repoURL: repoURL,
	}
}

func (p *GenericGitProvider) Name() string {
	return "git"
}

func (p *GenericGitProvider) Validate(r *http.Request, secret string) ([]byte, error) {
	// For generic git, we just read the body.
	// The caller can use X-Trigra-Secret if they want.
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (p *GenericGitProvider) ParsePushEvent(r *http.Request, payload []byte) (*providers.PushEvent, error) {
	// Generic git provider doesn't have a fixed webhook format.
	// Users can hit /webhook with anything to trigger a sync.
	return &providers.PushEvent{
		Owner:         "generic",
		Repo:          "repo",
		Ref:           "refs/heads/main",
		After:         "HEAD",
		ModifiedFiles: []string{"."}, // Always sync everything
	}, nil
}

func (p *GenericGitProvider) DownloadFile(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
	if p.repoURL == "" {
		return nil, fmt.Errorf("GIT_REPO_URL not configured")
	}

	repoDir := filepath.Join(p.baseDir, "repo")

	var r *git.Repository
	var err error

	if _, err = os.Stat(repoDir); os.IsNotExist(err) {
		r, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			URL: p.repoURL,
		})
	} else {
		r, err = git.PlainOpen(repoDir)
		if err == nil {
			w, _ := r.Worktree()
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		}
	}

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, err
	}

	if path == "." {
		// This is a special case: we don't return bytes for a directory here.
		// Higher level logic should handle this.
		return nil, fmt.Errorf("cannot download directory via DownloadFile")
	}

	return os.ReadFile(filepath.Join(repoDir, path))
}
