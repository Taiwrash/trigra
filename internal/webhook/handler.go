// Package webhook provides the HTTP handler for processing Git provider webhooks.
package webhook

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Taiwrash/trigra/internal/k8s"
	"github.com/Taiwrash/trigra/internal/manager"
	"github.com/Taiwrash/trigra/internal/providers"
)

// Handler handles Git provider webhook events using various providers.
type Handler struct {
	applier *k8s.Applier
	mgr     *manager.Manager
}

// NewHandler creates a new agnostic webhook handler.
func NewHandler(applier *k8s.Applier, mgr *manager.Manager) *Handler {
	return &Handler{
		applier: applier,
		mgr:     mgr,
	}
}

// ServeHTTP handles incoming webhook requests from any supported provider.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// 0. Identify project
	// URL format: /webhook/{project-name}
	projectName := strings.TrimPrefix(r.URL.Path, "/webhook/")
	if projectName == "" || projectName == "/webhook" {
		projectName = "default" // Fallback to singleton config
	}

	project := h.mgr.GetProject(projectName)
	if project == nil {
		log.Printf("ERROR: Project %s not found", projectName)
		http.Error(w, fmt.Sprintf("Project %s not found", projectName), http.StatusNotFound)
		return
	}

	log.Printf("INFO: Received webhook request for project: %s (Provider: %s)", projectName, project.Provider.Name())

	// 1. Validate webhook payload
	payload, err := project.Provider.Validate(r, project.WebhookSecret)
	if err != nil {
		log.Printf("ERROR: Invalid webhook payload for %s: %v", projectName, err)
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// 2. Parse push event
	event, err := project.Provider.ParsePushEvent(r, payload)
	if err != nil {
		log.Printf("WARNING: Generic parsing failed or non-push event: %v", err)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Event received but not processed (non-push)")
		return
	}

	// 3. Handle push event
	if err := h.handlePushEvent(ctx, project, event); err != nil {
		log.Printf("ERROR: Failed to handle push event: %v", err)
		http.Error(w, fmt.Sprintf("Failed to process push: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Successfully processed %s push event for project %s", project.Provider.Name(), projectName)
}

// handlePushEvent processes the generic push event.
func (h *Handler) handlePushEvent(ctx context.Context, p *manager.Project, event *providers.PushEvent) error {
	log.Printf("INFO: Processing push event from %s/%s on ref %s",
		event.Owner,
		event.Repo,
		event.Ref)

	// Filter for YAML files only
	yamlFiles := h.filterYAMLFiles(event.ModifiedFiles)

	if len(yamlFiles) == 0 {
		log.Printf("INFO: No YAML files found in push event")
		return nil
	}

	log.Printf("INFO: Found %d YAML file(s) to process: %s", len(yamlFiles), strings.Join(yamlFiles, ", "))

	// Process each YAML file
	for _, filename := range yamlFiles {
		if err := h.processFile(ctx, p, event, filename); err != nil {
			return fmt.Errorf("failed to process file %s: %w", filename, err)
		}
	}

	return nil
}

// processFile downloads and applies a single YAML file using the provider.
func (h *Handler) processFile(ctx context.Context, p *manager.Project, event *providers.PushEvent, filename string) error {
	log.Printf("INFO: Downloading file: %s", filename)

	// Download file content via provider
	content, err := p.Provider.DownloadFile(ctx, event.Owner, event.Repo, event.After, filename)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	log.Printf("INFO: Applying resources from file: %s to namespace %s", filename, p.TargetNamespace)

	// Apply the resource(s) to Kubernetes
	if err := h.applier.ApplyMultipleResources(ctx, content, p.TargetNamespace); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	log.Printf("SUCCESS: Applied resources from file: %s", filename)
	return nil
}

// filterYAMLFiles returns only YAML/YML files from the list.
func (h *Handler) filterYAMLFiles(files []string) []string {
	yamlFiles := []string{}
	for _, file := range files {
		if file == "." {
			return []string{"."} // Special case for full repo scan
		}
		if strings.HasSuffix(strings.ToLower(file), ".yaml") ||
			strings.HasSuffix(strings.ToLower(file), ".yml") {
			yamlFiles = append(yamlFiles, file)
		}
	}
	return yamlFiles
}
