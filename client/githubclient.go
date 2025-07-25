package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm/models"
	"time"
)

const githubAPI = "https://api.github.com"

func GetUserRepos(token string, since time.Time) ([]models.GithubRepo, error) {
	url := fmt.Sprintf("%s/user/repos?per_page=100", githubAPI)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var allRepos []models.GithubRepo
	if err := json.NewDecoder(resp.Body).Decode(&allRepos); err != nil {
		return nil, err
	}

	var filteredRepos []models.GithubRepo
	for _, repo := range allRepos {
		// Skip GitHub profile repo (e.g., yassir20191/yassir20191)
		if repo.Owner.Login == repo.Name {
			continue
		}
		updatedAt, err := time.Parse(time.RFC3339, repo.UpdatedAt)
		if err == nil && updatedAt.After(since) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

func GetUserMergedPRs(token, owner, repo string, since time.Time) ([]models.PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=closed&per_page=100", githubAPI, owner, repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var allPRs []models.PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&allPRs); err != nil {
		return nil, err
	}

	var filtered []models.PullRequest
	for _, pr := range allPRs {
		if pr.MergedAt == "" {
			continue
		}
		mergedAt, err := time.Parse(time.RFC3339, pr.MergedAt)
		if err != nil || mergedAt.Before(since) {
			continue
		}
		// Fetch detailed PR info
		detailURL := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", githubAPI, owner, repo, pr.Number)
		detailReq, _ := http.NewRequest("GET", detailURL, nil)
		detailReq.Header.Set("Authorization", "Bearer "+token)
		detailReq.Header.Set("Accept", "application/vnd.github+json")
		detailResp, err := client.Do(detailReq)
		if err != nil {
			continue
		}
		defer detailResp.Body.Close()

		var detailedPR models.PullRequest
		if err := json.NewDecoder(detailResp.Body).Decode(&detailedPR); err != nil {
			continue
		}
		filtered = append(filtered, detailedPR)
	}
	return filtered, nil
}

func GetRepoLanguages(token, owner, repo string) (map[string]int, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/languages", githubAPI, owner, repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var languages map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, err
	}
	return languages, nil
}

func GetGitHubUsername(token string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}
	return user.Login, nil
}

func GetUserCommits(token, owner, repo, username string, since time.Time) (int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?author=%s&since=%s", owner, repo, username, since.Format(time.RFC3339))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var commits []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return 0, err
	}

	return len(commits), nil
}
