// Package main is the entry point for the Trigra GitOps controller.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Taiwrash/trigra/internal/config"
	"github.com/Taiwrash/trigra/internal/k8s"
	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/Taiwrash/trigra/internal/providers/bitbucket"
	"github.com/Taiwrash/trigra/internal/providers/git"
	"github.com/Taiwrash/trigra/internal/providers/gitea"
	"github.com/Taiwrash/trigra/internal/providers/github"
	"github.com/Taiwrash/trigra/internal/providers/gitlab"
	"github.com/Taiwrash/trigra/internal/webhook"
)

var (
	// Version is the application version, set at build time.
	Version = "dev"
	// BuildTime is the time the application was built, set at build time.
	BuildTime = "unknown"
)

func main() {
	log.Printf("Starting Trigra GitOps Controller Version=%s BuildTime=%s", Version, BuildTime)

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded: InCluster=%v, Namespace=%s, Provider=%s", cfg.InCluster, cfg.Namespace, cfg.GitProvider)

	// 2. Initialize Kubernetes Applier
	applier, err := k8s.NewApplier(cfg.InCluster)
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes applier: %v", err)
	}

	// 3. Initialize Git Provider
	var provider providers.Provider
	switch cfg.GitProvider {
	case "github":
		provider = github.NewProvider(cfg.GitToken)
	case "gitlab":
		provider = gitlab.NewProvider(cfg.GitBaseURL, cfg.GitToken)
	case "gitea":
		provider = gitea.NewProvider(cfg.GitBaseURL, cfg.GitToken)
	case "bitbucket":
		provider = bitbucket.NewProvider(os.Getenv("BITBUCKET_USER"), cfg.GitToken)
	case "git":
		provider = git.NewProvider(cfg.GitRepoURL, cfg.GitSSHKeyFile)
	default:
		log.Fatalf("Unsupported git provider: %s", cfg.GitProvider)
	}

	// 4. Automated Webhook Setup
	if cfg.PublicURL != "" {
		setupWebhooks(cfg, provider)
	}

	// 5. Initialize Webhook Handler
	handler := webhook.NewHandler(applier, provider, cfg.WebhookSecret, cfg.Namespace)

	// 6. Start Server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:           nil,             // Use DefaultServeMux
		ReadHeaderTimeout: 3 * time.Second, // Fix G112: Slowloris Attack
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Register webhook endpoint
	http.Handle("/webhook", handler)

	// Add health check endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Ready")
	})

	// Create a channel to listen for OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		log.Printf("Server listening on port %d", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for signal
	sig := <-sigs
	log.Printf("Received signal %s, shutting down...", sig)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func setupWebhooks(cfg *config.Config, p providers.Provider) {
	owner := cfg.GitOwner
	repo := cfg.GitRepo

	// Try to auto-parse if missing
	if (owner == "" || repo == "") && cfg.GitRepoURL != "" {
		// Handle HTTPS URLs
		if strings.HasPrefix(cfg.GitRepoURL, "http") {
			path := strings.TrimPrefix(cfg.GitRepoURL, "https://")
			path = strings.TrimPrefix(path, "http://")
			path = strings.TrimSuffix(path, ".git")
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				repo = parts[len(parts)-1]
				owner = strings.Join(parts[1:len(parts)-1], "/")
			}
		} else if strings.Contains(cfg.GitRepoURL, "@") && strings.Contains(cfg.GitRepoURL, ":") {
			// Handle SSH URLs like git@github.com:owner/repo.git or git@gitlab.com:group/subgroup/repo.git
			path := strings.Split(cfg.GitRepoURL, ":")[1]
			path = strings.TrimSuffix(path, ".git")
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				repo = parts[len(parts)-1]
				owner = strings.Join(parts[:len(parts)-1], "/")
			}
		}
	}

	if owner == "" || repo == "" {
		log.Printf("WARNING: Automated webhook setup skipped - could not determine owner/repo")
		return
	}

	webhookURL := fmt.Sprintf("%s/webhook", strings.TrimSuffix(cfg.PublicURL, "/"))
	log.Printf("INFO: Attempting to configure webhook for %s/%s at %s", owner, repo, webhookURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := p.SetupWebhook(ctx, owner, repo, webhookURL, cfg.WebhookSecret); err != nil {
		log.Printf("WARNING: Automated webhook setup failed: %v", err)
	} else {
		log.Printf("SUCCESS: Webhook ensured for %s/%s", owner, repo)
	}
}
