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

// Provider implements providers.Provider by cloning repositories locally.
type Provider struct {
	baseDir string
	repoURL string
}

// NewProvider creates a new generic Git provider.
func NewProvider(repoURL string) *Provider {
	tmpDir, _ := os.MkdirTemp("", "trigra-git-*")
	return &Provider{
		baseDir: tmpDir,
		repoURL: repoURL,
	}
}

// Name returns "git".
func (p *Provider) Name() string {
	return "git"
}

// Validate handles generic webhook validation.
func (p *Provider) Validate(r *http.Request, _ string) ([]byte, error) {
	// For generic git, we just read the body.
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// ParsePushEvent parses a generic push event.
func (p *Provider) ParsePushEvent(_ *http.Request, _ []byte) (*providers.PushEvent, error) {
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
func (p *Provider) DownloadFile(_ context.Context, _, _, ref, path string) ([]byte, error) {
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

	// We use ReadFile on a sanitized and validated path.
	// #nosec G304
	return os.ReadFile(cleanPath)
}
