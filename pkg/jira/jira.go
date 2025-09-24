package jira

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-jira"
)

// Client represents a client for interacting with the Jira API.
type Client struct {
	client *jira.Client
}

// NewClient creates a new instance of the Jira client using an email and API token.
func NewClient(baseURL, email, token string) (*Client, error) {
	tp := jira.BasicAuthTransport{
		Username: email,
		Password: token,
	}

	client, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	return &Client{client: client}, nil
}

// Issue represents a simplified structure for a Jira issue.
type Issue struct {
	Key      string
	Assignee string
	Status   string
	Project  string
}

// GetIssues retrieves issues for a given JQL query.
func (c *Client) GetIssues(ctx context.Context, jql string) ([]Issue, error) {
	// Execute a JQL query against the Jira API
	jiraIssues, _, err := c.client.Issue.SearchWithContext(ctx, jql, nil)
	if err != nil {
		return nil, fmt.Errorf("error executing JQL query '%s': %w", jql, err)
	}

	var issues []Issue
	for _, issue := range jiraIssues {
		assignee := "unassigned"
		if issue.Fields.Assignee != nil {
			assignee = issue.Fields.Assignee.Name
		}

		proj := ""
		if issue.Fields.Project.Name != "" {
			proj = issue.Fields.Project.Name
		}

		issues = append(issues, Issue{
			Key:      issue.Key,
			Assignee: assignee,
			Status:   issue.Fields.Status.Name,
			Project:  proj,
		})
	}

	return issues, nil
}

