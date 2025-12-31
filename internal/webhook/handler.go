package webhook

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Taiwrash/trigra/internal/k8s"
	"github.com/google/go-github/v79/github"
)

// Handler handles GitHub webhook events
type Handler struct {
	applier       *k8s.Applier
	githubClient  *github.Client
	webhookSecret string
	namespace     string
}

// NewHandler creates a new webhook handler
func NewHandler(applier *k8s.Applier, githubClient *github.Client, webhookSecret, namespace string) *Handler {
	return &Handler{
		applier:       applier,
		githubClient:  githubClient,
		webhookSecret: webhookSecret,
		namespace:     namespace,
	}
}

// ServeHTTP handles incoming webhook requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Validate webhook payload
	payload, err := github.ValidatePayload(r, []byte(h.webhookSecret))
	if err != nil {
		log.Printf("ERROR: Invalid webhook payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Parse webhook event
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("ERROR: Failed to parse webhook: %v", err)
		http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
		return
	}

	// Handle different event types
	switch e := event.(type) {
	case *github.PingEvent:
		log.Printf("INFO: Received ping event from webhook")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "pong")
		return

	case *github.PushEvent:
		if err := h.handlePushEvent(ctx, e); err != nil {
			log.Printf("ERROR: Failed to handle push event: %v", err)
			http.Error(w, fmt.Sprintf("Failed to process push: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Successfully processed push event")
		return

	default:
		log.Printf("WARNING: Unhandled webhook event type: %T", event)
		http.Error(w, "Event type not supported", http.StatusNotImplemented)
		return
	}
}

// handlePushEvent processes a GitHub push event
func (h *Handler) handlePushEvent(ctx context.Context, event *github.PushEvent) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()

	log.Printf("INFO: Processing push event from %s/%s", owner, repo)

	var files []string
	if event.GetBefore() == "0000000000000000000000000000000000000000" {
		// New branch/push - get files from the single commit
		commit, _, err := h.githubClient.Repositories.GetCommit(ctx, owner, repo, event.GetAfter(), nil)
		if err != nil {
			return fmt.Errorf("failed to get commit: %w", err)
		}
		for _, f := range commit.Files {
			files = append(files, f.GetFilename())
		}
	} else {
		// Existing branch - compare range
		comp, _, err := h.githubClient.Repositories.CompareCommits(ctx, owner, repo, event.GetBefore(), event.GetAfter(), nil)
		if err != nil {
			return fmt.Errorf("failed to compare commits: %w", err)
		}
		for _, f := range comp.Files {
			files = append(files, f.GetFilename())
		}
	}

	// Filter for YAML files only
	yamlFiles := h.filterYAMLFiles(files)

	if len(yamlFiles) == 0 {
		log.Printf("INFO: No YAML files found in push event")
		return nil
	}

	log.Printf("INFO: Found %d YAML file(s) to process: %s", len(yamlFiles), strings.Join(yamlFiles, ", "))

	// Process each YAML file
	for _, filename := range yamlFiles {
		if err := h.processFile(ctx, event, filename); err != nil {
			return fmt.Errorf("failed to process file %s: %w", filename, err)
		}
	}

	return nil
}

// processFile downloads and applies a single YAML file
func (h *Handler) processFile(ctx context.Context, event *github.PushEvent, filename string) error {
	log.Printf("INFO: Downloading file: %s", filename)

	// Download file content from GitHub
	fileReader, _, err := h.githubClient.Repositories.DownloadContents(
		ctx,
		event.Repo.GetOwner().GetName(),
		event.Repo.GetName(),
		filename,
		&github.RepositoryContentGetOptions{
			Ref: event.GetAfter(), // Use the commit SHA
		},
	)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer fileReader.Close()

	// Read file content
	content, err := io.ReadAll(fileReader)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	log.Printf("INFO: Applying resources from file: %s", filename)

	// Apply the resource(s) to Kubernetes
	if err := h.applier.ApplyMultipleResources(ctx, content, h.namespace); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	log.Printf("SUCCESS: Applied resources from file: %s", filename)
	return nil
}

// filterYAMLFiles returns only YAML/YML files from the list
func (h *Handler) filterYAMLFiles(files []string) []string {
	yamlFiles := []string{}

	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".yaml") ||
			strings.HasSuffix(strings.ToLower(file), ".yml") {
			yamlFiles = append(yamlFiles, file)
		}
	}

	return yamlFiles
}
