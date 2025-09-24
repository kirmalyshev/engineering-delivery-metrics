package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v75/github"
)

// Client is a client for interacting with the GitHub API.
type Client struct {
	client *github.Client
	owner  string
	repo   string
}

// NewClient creates a new instance of the GitHub client.
// token is the Personal Access Token for authentication.
// owner is the repository owner (organization or user).
// repo is the name of the repository.
func NewClient(token, owner, repo string) *Client {
	// For authenticated requests, an http.Client with a token is used
	println("Token: `%s`", token)
	client := github.NewClient(nil).WithAuthToken(token)

	return &Client{
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

// CommitInfo represents summarized information about a commit.
type CommitInfo struct {
	Author string
	SHA    string
	Date   string
}

// PullRequestInfo represents summarized information about a Pull Request.
type PullRequestInfo struct {
	Author string
	Title  string
	Number int
	State  string
}
// RepoInfo represents summarized information about a repository.
type RepoInfo struct {
	Owner string
	Name  string
}

// GetRepos fetches all repositories accessible by the authenticated user.
func (c *Client) GetRepos(ctx context.Context) ([]RepoInfo, error) {
	// Passing "" as the user will list repositories for the authenticated user.
	// We'll use pagination to fetch all repositories, not just the first page.
	var allRepos []*github.Repository

	//get all repos for org
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := c.client.Repositories.ListByOrg(ctx, c.owner, opt)
		if err != nil {
			return nil, fmt.Errorf("error getting repositories from GitHub: %w", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var repoInfos []RepoInfo
	for _, repo := range allRepos {
		if repo.GetOwner() != nil {
			info := RepoInfo{
				Owner: repo.GetOwner().GetLogin(),
				Name:  repo.GetName(),
			}
			repoInfos = append(repoInfos, info)
		}
	}

	return repoInfos, nil
}


// GetCommits retrieves the list of commits for the repository.
func (c *Client) GetCommits(ctx context.Context) ([]CommitInfo, error) {
	commits, _, err := c.client.Repositories.ListCommits(ctx, c.owner, c.repo, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting commits from GitHub: %w", err)
	}

	var commitInfos []CommitInfo
	for _, commit := range commits {
		author := "unknown"
		// The commit author may not be associated with a GitHub user
		if commit.GetAuthor() != nil {
			author = commit.GetAuthor().GetLogin()
		}

		info := CommitInfo{
			Author: author,
			SHA:    commit.GetSHA(),
			Date:   commit.GetCommit().GetAuthor().GetDate().String(),
		}
		commitInfos = append(commitInfos, info)
	}

	return commitInfos, nil
}

// GetPullRequests retrieves the list of Pull Requests for the repository.
func (c *Client) GetPullRequests(ctx context.Context) ([]PullRequestInfo, error) {
	// Options to get all PRs (open and closed)
	opts := &github.PullRequestListOptions{
		State: "all",
	}
	prs, _, err := c.client.PullRequests.List(ctx, c.owner, c.repo, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting pull requests: %w", err)
	}

	var prInfos []PullRequestInfo
	for _, pr := range prs {
		author := "unknown"
		if pr.GetUser() != nil {
			author = pr.GetUser().GetLogin()
		}
		info := PullRequestInfo{
			Author: author,
			Title:  pr.GetTitle(),
			Number: pr.GetNumber(),
			State:  pr.GetState(),
		}
		prInfos = append(prInfos, info)
	}

	return prInfos, nil
}

