package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pm/models"
	"time"
)

const githubAPI = "https://api.github.com"

func getUserRepos(token string, since time.Time) ([]models.GithubRepo, error) {
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
		updatedAt, err := time.Parse(time.RFC3339, repo.UpdatedAt)
		if err == nil && updatedAt.After(since) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <YYYY-MM-DD>")
	}
	sinceDate, err := time.Parse("2006-01-02", os.Args[1])
	if err != nil {
		log.Fatalf("Invalid date format: %v", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Set GITHUB_TOKEN environment variable.")
	}

	repos, err := getUserRepos(token, sinceDate)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Repositories updated since %s:\n", sinceDate.Format("2006-01-02"))
	for _, repo := range repos {
		fmt.Println(" -", repo.FullName)
	}
}

const test = "test"
