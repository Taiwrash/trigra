// Package webhook provides the HTTP handler for processing Git provider webhooks.
package webhook

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Taiwrash/trigra/internal/k8s"
	"github.com/Taiwrash/trigra/internal/providers"
)

// Handler handles Git provider webhook events using various providers.
type Handler struct {
	applier       *k8s.Applier
	provider      providers.Provider
	webhookSecret string
	namespace     string
}

// NewHandler creates a new agnostic webhook handler.
func NewHandler(applier *k8s.Applier, provider providers.Provider, webhookSecret, namespace string) *Handler {
	return &Handler{
		applier:       applier,
		provider:      provider,
		webhookSecret: webhookSecret,
		namespace:     namespace,
	}
}

// ServeHTTP handles incoming webhook requests from any supported provider.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	log.Printf("INFO: Received webhook request for provider: %s", h.provider.Name())

	// 1. Validate webhook payload
	payload, err := h.provider.Validate(r, h.webhookSecret)
	if err != nil {
		log.Printf("ERROR: Invalid webhook payload for %s: %v", h.provider.Name(), err)
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// 2. Parse push event
	event, err := h.provider.ParsePushEvent(r, payload)
	if err != nil {
		// Some events like 'ping' might not be push events. Check headers if needed,
		// but for now we follow the generic push event flow.
		log.Printf("WARNING: Generic parsing failed or non-push event: %v", err)
		// We can return 200 for health checks/pings that aren't push events
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Event received but not processed (non-push)")
		return
	}

	// 3. Handle push event
	if err := h.handlePushEvent(ctx, event); err != nil {
		log.Printf("ERROR: Failed to handle push event: %v", err)
		http.Error(w, fmt.Sprintf("Failed to process push: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Successfully processed %s push event", h.provider.Name())
}

// handlePushEvent processes the generic push event.
func (h *Handler) handlePushEvent(ctx context.Context, event *providers.PushEvent) error {
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
		if err := h.processFile(ctx, event, filename); err != nil {
			return fmt.Errorf("failed to process file %s: %w", filename, err)
		}
	}

	return nil
}

// processFile downloads and applies a single YAML file using the provider.
func (h *Handler) processFile(ctx context.Context, event *providers.PushEvent, filename string) error {
	log.Printf("INFO: Downloading file: %s", filename)

	// Download file content via provider
	content, err := h.provider.DownloadFile(ctx, event.Owner, event.Repo, event.After, filename)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	log.Printf("INFO: Applying resources from file: %s", filename)

	// Apply the resource(s) to Kubernetes
	if err := h.applier.ApplyMultipleResources(ctx, content, h.namespace); err != nil {
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
