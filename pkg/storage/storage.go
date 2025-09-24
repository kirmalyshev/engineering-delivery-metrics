package storage

import (
	"log"
	"net/http"
	"sync"

	"delivery-dashboard/pkg/github"
	"delivery-dashboard/pkg/jira"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsStorage is responsible for storing metrics and exposing them to Prometheus.
type MetricsStorage struct {
	// We use sync.Mutex for safe concurrent access to metrics from different goroutines
	mu sync.Mutex

	// Define Prometheus metrics
	githubCommitsTotal      *prometheus.CounterVec
	githubPullRequestsTotal *prometheus.GaugeVec
	jiraIssues              *prometheus.GaugeVec
}

// NewMetricsStorage creates and registers metrics with Prometheus.
func NewMetricsStorage() *MetricsStorage {
	s := &MetricsStorage{
		githubCommitsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "github_commits_total",
				Help: "Total number of commits by author and repository.",
			},
			[]string{"author", "repo"}, // Labels for the metric
		),
		githubPullRequestsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "github_pull_requests_total",
				Help: "Number of Pull Requests by status.",
			},
			[]string{"repo", "state", "author"},
		),
		jiraIssues: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jira_issues_status_count",
				Help: "Number of Jira issues in a specific status.",
			},
			[]string{"project", "status", "assignee"},
		),
	}

	// Register metrics
	prometheus.MustRegister(s.githubCommitsTotal)
	prometheus.MustRegister(s.githubPullRequestsTotal)
	prometheus.MustRegister(s.jiraIssues)

	return s
}

// UpdateGithubCommits updates metrics based on commit data.
func (s *MetricsStorage) UpdateGithubCommits(commits []github.CommitInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// This logic is still simplified. In a production system,
	// you would need to store the SHA of the last processed commit
	// to avoid counting the same commits repeatedly.
	for _, commit := range commits {
		// TODO: Get repo name from config
		s.githubCommitsTotal.WithLabelValues(commit.Author, "your-repo").Inc()
	}
}

// UpdateGithubPullRequests updates the metric for Pull Requests.
func (s *MetricsStorage) UpdateGithubPullRequests(prs []github.PullRequestInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Reset old values before updating
	s.githubPullRequestsTotal.Reset()
	for _, pr := range prs {
		// TODO: Get repo name from config
		s.githubPullRequestsTotal.WithLabelValues("your-repo", pr.State, pr.Author).Inc()
	}
}

// UpdateJiraIssues updates metrics based on issue data.
func (s *MetricsStorage) UpdateJiraIssues(issues []jira.Issue) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Reset previous values to avoid stale data
	s.jiraIssues.Reset()
	for _, issue := range issues {
		s.jiraIssues.WithLabelValues(issue.Project, issue.Status, issue.Assignee).Inc()
	}
}

// StartServer starts an HTTP server for the /metrics endpoint.
func (s *MetricsStorage) StartServer(addr string) {
	log.Printf("Metrics server listening on: %s/metrics", addr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}

