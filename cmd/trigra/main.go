package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Taiwrash/trigra/internal/config"
	"github.com/Taiwrash/trigra/internal/k8s"
	"github.com/Taiwrash/trigra/internal/webhook"
	"github.com/google/go-github/v79/github"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	log.Printf("Starting Trigra GitOps Controller Version=%s BuildTime=%s", Version, BuildTime)

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded: InCluster=%v, Namespace=%s", cfg.InCluster, cfg.Namespace)

	// 2. Initialize Kubernetes Applier
	applier, err := k8s.NewApplier(cfg.InCluster)
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes applier: %v", err)
	}

	// 3. Initialize GitHub Client
	var client *github.Client
	if cfg.GitHubToken != "" {
		client = github.NewClient(nil).WithAuthToken(cfg.GitHubToken)
	} else {
		client = github.NewClient(nil)
	}

	// 4. Initialize Webhook Handler
	handler := webhook.NewHandler(applier, client, cfg.WebhookSecret, cfg.Namespace)

	// 5. Start Server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:           nil, // Use DefaultServeMux
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Register webhook endpoint
	http.Handle("/webhook", handler)

	// Add health check endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
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
