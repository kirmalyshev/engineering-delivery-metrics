package collector

import (
	"context"
	"log"
	"time"

	"delivery-dashboard/pkg/github"
	"delivery-dashboard/pkg/jira"
	"delivery-dashboard/pkg/storage"
)

// Collector manages data collection from various sources.
type Collector struct {
	githubClient *github.Client
	jiraClient   *jira.Client
	storage      *storage.MetricsStorage
}

// New creates a new instance of the Collector.
func New(gh *github.Client, jr *jira.Client, st *storage.MetricsStorage) *Collector {
	return &Collector{
		githubClient: gh,
		jiraClient:   jr,
		storage:      st,
	}
}

// Run starts the periodic data collection.
func (c *Collector) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("Collector started. Initial data collection...")
	c.collect(ctx) // First collection runs immediately on start

	for {
		select {
		case <-ticker.C:
			log.Println("Starting scheduled data collection...")
			c.collect(ctx)
		case <-ctx.Done():
			log.Println("Collector stopped.")
			return
		}
	}
}

// collect performs a single data collection cycle.
func (c *Collector) collect(ctx context.Context) {
	// --- Collect data from GitHub ---
	if c.githubClient == nil {
		log.Fatalf("GitHub client is not configured")
	}
	// Collect commits
	//commits, err := c.githubClient.GetCommits(ctx)
	//if err != nil {
	//	log.Printf("Error collecting commits from GitHub: %v\n", err)
	//} else {
	//	log.Printf("Fetched %d commits from GitHub\n", len(commits))
	//	c.storage.UpdateGithubCommits(commits)
	//}

	// Collect Pull Requests
	repos, err := c.githubClient.GetRepos(ctx)
	if err != nil {
		log.Printf("[github] Error collecting Repos: %v\n", err)
	} else {
		log.Printf("[github] Fetched %d repos\n", len(repos))
		log.Printf("[github] repos: %+v", repos)
	}

	// Collect Pull Requests
	prs, err := c.githubClient.GetPullRequests(ctx)
	if err != nil {
		log.Printf("[github] Error collecting Pull Requests: %v\n", err)
	} else {
		log.Printf("[github] Fetched %d Pull Requests\n", len(prs))
		c.storage.UpdateGithubPullRequests(prs)
	}

	// --- Collect data from Jira ---
	//if c.jiraClient != nil {
	//	// As an example, use a JQL query for all issues in project "PROJ"
	//	// TODO: Move this to configuration
	//	jql := "project = GT"
	//	issues, err := c.jiraClient.GetIssues(ctx, jql)
	//	if err != nil {
	//		log.Printf("[jira] Error collecting issues: %v\n", err)
	//	} else {
	//		log.Printf("[jira] Fetched %d issues\n", len(issues))
	//		c.storage.UpdateJiraIssues(issues)
	//	}
	//}

	log.Println("Data collection finished.")
}

