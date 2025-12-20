// Package git provides a generic Git implementation using local cloning.
package git

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GenericGitProvider implements providers.Provider by cloning repositories locally.
type GenericGitProvider struct {
	baseDir string
	repoURL string
}

// NewGenericGitProvider creates a new generic Git provider.
func NewGenericGitProvider(repoURL string) *GenericGitProvider {
	tmpDir, _ := os.MkdirTemp("", "trigra-git-*")
	return &GenericGitProvider{
		baseDir: tmpDir,
		repoURL: repoURL,
	}
}

// Name returns "git".
func (p *GenericGitProvider) Name() string {
	return "git"
}

// Validate handles generic webhook validation.
func (p *GenericGitProvider) Validate(r *http.Request, _ string) ([]byte, error) {
	// For generic git, we just read the body.
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// ParsePushEvent parses a generic push event.
func (p *GenericGitProvider) ParsePushEvent(_ *http.Request, _ []byte) (*providers.PushEvent, error) {
	// Generic git provider doesn't have a fixed webhook format.
	return &providers.PushEvent{
		Owner:         "generic",
		Repo:          "repo",
		Ref:           "refs/heads/main",
		After:         "HEAD",
		ModifiedFiles: []string{"."}, // Always sync everything
	}, nil
}

// DownloadFile "downloads" a file by reading it from the local clone.
func (p *GenericGitProvider) DownloadFile(_ context.Context, _, _, ref, path string) ([]byte, error) {
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
			// Pull to ensure we have the latest. Ignore already up to date.
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err == git.NoErrAlreadyUpToDate {
				err = nil
			}
		}
	}

	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	if ref != "HEAD" {
		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(ref),
		})
	} else {
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName("main"), // Fallback/default logic
		})
	}

	if err != nil {
		return nil, err
	}

	// Sanitize path to prevent path traversal (Fix G304)
	cleanPath := filepath.Join(repoDir, filepath.Clean(path))
	if !strings.HasPrefix(cleanPath, filepath.Clean(repoDir)) {
		return nil, fmt.Errorf("traversal attack detected: %s", path)
	}

	return os.ReadFile(cleanPath)
}
