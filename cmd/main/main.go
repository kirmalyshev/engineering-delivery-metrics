package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"engineering-delivery-metrics/pkg/collector"
	"engineering-delivery-metrics/pkg/github"
	"engineering-delivery-metrics/pkg/jira"
	"engineering-delivery-metrics/pkg/storage"

	"github.com/joho/godotenv" // Import the library for .env files
)

func main() {
	// Load environment variables from the .env file in the project root
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING: .env file not found. Using system environment variables.")
	}

	// --- Configuration from environment variables ---
	githubToken := os.Getenv("GITHUB_TOKEN")
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")

	jiraURL := os.Getenv("JIRA_URL")
	jiraEmail := os.Getenv("JIRA_EMAIL")
	jiraApiToken := os.Getenv("JIRA_API_TOKEN")

	// Check for required variables
	if githubToken == "" || githubOwner == "" || githubRepo == "" {
		log.Fatal("GitHub variables (GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO) must be set.")
	}
	if jiraURL == "" || jiraEmail == "" || jiraApiToken == "" {
		log.Fatal("Jira variables (JIRA_URL, JIRA_EMAIL, JIRA_API_TOKEN) must be set.")
	}

	intervalStr := os.Getenv("COLLECTION_INTERVAL")
	if intervalStr == "" {
		intervalStr = "1m" // Default value
	}
	println("Collection interval: %d", intervalStr)
	collectionInterval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("Invalid format for collection interval: %v", err)
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":9091" // Default value
	}

	// --- Component Initialization ---

	// 1. API Clients
	ghClient := github.NewClient(githubToken, githubOwner, githubRepo)

	jiraClient, err := jira.NewClient(jiraURL, jiraEmail, jiraApiToken)
	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}

	// 2. Metrics Storage
	metricsStorage := storage.NewMetricsStorage()

	// 3. Collector
	dataCollector := collector.New(ghClient, jiraClient, metricsStorage)

	// --- Application Start ---

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the metrics server in a separate goroutine
	go metricsStorage.StartServer(serverAddr)

	// Start the data collector in a separate goroutine
	go dataCollector.Run(ctx, collectionInterval)

	log.Println("Application started. Press Ctrl+C to exit.")

	// Wait for a signal to terminate
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received, stopping application...")
	cancel() // Cancel the context to stop the goroutines

	// Allow time for graceful shutdown
	time.Sleep(2 * time.Second)
	log.Println("Application stopped.")
}

