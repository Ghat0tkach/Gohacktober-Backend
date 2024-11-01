package github

import (
	"context"
	"fmt"
	"time"
	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
	"github.com/Ghat0tkach/Gohacktober-Backend/config"
	"sync"
)

var client *github.Client


func Init(cfg *config.Config) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	fmt.Println("GitHub client initialized")
}

func FetchHacktoberfestContributions(org, username string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	fmt.Printf("Fetching Hacktoberfest contributions for org: %s and user: %s\n", org, username)

	repos, err := fetchOrgHacktoberfestRepos(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("error fetching Hacktoberfest repos: %v", err)
	}
	fmt.Printf("Found %d Hacktoberfest repos for org %s\n", len(repos), org)

	var contributions []map[string]interface{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, repo := range repos {
		wg.Add(1)
		go func(repo *github.Repository) {
			defer wg.Done()
			fmt.Printf("Fetching contributions for repo: %s\n", repo.GetName())
			repoContributions, err := fetchUserContributions(ctx, org, repo.GetName(), username)
			if err != nil {
				fmt.Printf("Error fetching contributions for %s: %v\n", repo.GetName(), err)
				return
			}
			
			repoData := map[string]interface{}{
				"repo_name":     repo.GetName(),
				"contributions": repoContributions,
			}
			
			mu.Lock()
			contributions = append(contributions, repoData)
			mu.Unlock()
			fmt.Printf("Finished fetching contributions for repo: %s\n", repo.GetName())
		}(repo)
	}

	wg.Wait()
	fmt.Println("Finished fetching all contributions")

	return contributions, nil
}
func fetchOrgHacktoberfestRepos(ctx context.Context, account string) ([]*github.Repository, error) {
	fmt.Printf("Fetching Hacktoberfest repos for account: %s\n", account)
	var allRepos []*github.Repository
	var err error
    startDate := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2024, time.November, 1, 0, 0, 0, 0, time.UTC)
	// First, try to fetch as an organization
	orgRepos, err := fetchReposAsOrg(ctx, account)
	if err == nil {
		allRepos = orgRepos
	} else {
		fmt.Printf("Failed to fetch repos as organization, trying as user: %v\n", err)
		// If fetching as an organization fails, try as a user
		userRepos, err := fetchReposAsUser(ctx, account)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repos as both org and user: %v", err)
		}
		allRepos = userRepos
	}

	var hacktoberfestRepos []*github.Repository
	for _, repo := range allRepos {
		topics, _, err := client.Repositories.ListAllTopics(ctx, account, repo.GetName())
		if err != nil {
			return nil, err
		}
		
		for _, topic := range topics {
			if topic == "hacktoberfest" {
				// Check if the repo creation date is within the specified range
				if repo.GetCreatedAt().After(startDate) && repo.GetCreatedAt().Before(endDate) {
					hacktoberfestRepos = append(hacktoberfestRepos, repo)
					fmt.Printf("Found Hacktoberfest repo: %s, Created At: %s\n", repo.GetName(), repo.GetCreatedAt().Format("2006-01-02"))
				}
				break
			}
		}
	}

	return hacktoberfestRepos, nil
}

func fetchReposAsOrg(ctx context.Context, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func fetchReposAsUser(ctx context.Context, username string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, username, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func fetchUserContributions(ctx context.Context, owner, repo, username string) (map[string]interface{}, error) {
	fmt.Printf("Fetching contributions for user %s in repo %s/%s\n", username, owner, repo)
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		Creator: username,
		State:   "all",
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching issues: %v", err)
	}
	fmt.Printf("Found %d issues for user %s in repo %s/%s\n", len(issues), username, owner, repo)

	pulls, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "all",
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching pull requests: %v", err)
	}

	userPulls := filterPullRequestsByUser(pulls, username)
	fmt.Printf("Found %d pull requests for user %s in repo %s/%s\n", len(userPulls), username, owner, repo)

	return map[string]interface{}{
		"issues": summarizeIssues(issues),
		"pulls":  summarizePullRequests(userPulls),
	}, nil
}

func filterPullRequestsByUser(pulls []*github.PullRequest, username string) []*github.PullRequest {
	var userPulls []*github.PullRequest
	for _, pull := range pulls {
		if pull.User.GetLogin() == username {
			userPulls = append(userPulls, pull)
		}
	}
	return userPulls
}

func summarizeIssues(issues []*github.Issue) []map[string]interface{} {
	var summary []map[string]interface{}
	for _, issue := range issues {
		summary = append(summary, map[string]interface{}{
			"number":     issue.GetNumber(),
			"title":      issue.GetTitle(),
			"state":      issue.GetState(),
			"created_at": issue.GetCreatedAt(),
			"updated_at": issue.GetUpdatedAt(),
			"html_url":   issue.GetHTMLURL(),
		})
	}
	return summary
}

func summarizePullRequests(pulls []*github.PullRequest) []map[string]interface{} {
	var summary []map[string]interface{}
	for _, pull := range pulls {
		summary = append(summary, map[string]interface{}{
			"number":     pull.GetNumber(),
			"title":      pull.GetTitle(),
			"state":      pull.GetState(),
			"created_at": pull.GetCreatedAt(),
			"updated_at": pull.GetUpdatedAt(),
			"html_url":   pull.GetHTMLURL(),
		})
	}
	return summary
}