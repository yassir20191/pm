package service

import (
	"fmt"
	"github.com/google/go-github/v55/github"
	"log"
	githubclient "pm/client"
	"pm/models"
	"strings"
	"time"
)

func generateRepoLevelMetrics(token, owner, repoName string, sinceDate time.Time) string {
	var b strings.Builder

	// Languages
	languages, err := githubclient.GetRepoLanguages(token, owner, repoName)
	if err != nil {
		b.WriteString(fmt.Sprintf("   ⚠️ Error fetching languages: %v\n", err))
	} else if len(languages) > 0 {
		var langList []string
		for lang := range languages {
			langList = append(langList, lang)
		}
		b.WriteString(fmt.Sprintf("   Languages: %s\n", strings.Join(langList, ", ")))
	}

	// Pull Requests
	prs, err := githubclient.GetUserMergedPRs(token, owner, repoName, sinceDate)
	if err != nil {
		b.WriteString(fmt.Sprintf("   ⚠️ Error fetching PRs: %v\n", err))
		return b.String()
	}

	for _, pr := range prs {
		b.WriteString(fmt.Sprintf("   🟢 PR: %s\n", pr.Title))
		b.WriteString(fmt.Sprintf("     Description : %s\n", pr.Body))
		b.WriteString(fmt.Sprintf("     📁 Files changed: %d\n", pr.ChangedFiles))
		b.WriteString(fmt.Sprintf("     ✍️ Lines changed: +%d -%d\n", pr.Additions, pr.Deletions))
	}

	return b.String()
}

func generateIssueEngagementMetrics(repos []models.GithubRepo) string {
	issueCount := 0
	for _, repo := range repos {
		issueCount += repo.IssueFixCount
	}
	if issueCount == 0 {
		return ""
	}
	return fmt.Sprintf("🐞 Issues Fixed: %d", issueCount)
}

func generateCollaborationMetrics(token string, repos []models.GithubRepo, since time.Time) string {
	reviewCount := 0
	for _, repo := range repos {
		reviewCount += repo.ReviewCount
	}
	if reviewCount == 0 {
		return ""
	}
	return fmt.Sprintf("🔍 PRs Reviewed: %d", reviewCount)
}

func generateCommitLevelMetrics(token string, repos []models.GithubRepo, sincetime time.Time) string {
	username, err := githubclient.GetGitHubUsername(token)
	if err != nil {
		log.Println("⚠️ Could not retrieve GitHub username:", err)
		return ""
	}

	commitCount := 0
	for _, repo := range repos {
		parts := strings.Split(repo.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repoName := parts[0], parts[1]

		count, err := githubclient.GetUserCommits(token, owner, repoName, username, sincetime)
		if err == nil {
			commitCount += count
		}
	}
	if commitCount == 0 {
		return ""
	}
	return fmt.Sprintf("🔢 Total Commits: %d", commitCount)
}

func generatePullRequestMetrics(token string, repos []models.GithubRepo, since time.Time) string {
	totalPRs := 0
	var totalMergeTime time.Duration
	prCount := 0

	for _, repo := range repos {
		parts := strings.Split(repo.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repoName := parts[0], parts[1]
		prs, err := githubclient.GetUserMergedPRs(token, owner, repoName, since)
		if err != nil {
			continue
		}
		for _, pr := range prs {
			totalPRs++
			createdAt, err1 := time.Parse(time.RFC3339, pr.CreatedAt)
			mergedAt, err2 := time.Parse(time.RFC3339, pr.MergedAt)
			if err1 == nil && err2 == nil {
				totalMergeTime += mergedAt.Sub(createdAt)
				prCount++
			}
		}
	}

	if totalPRs == 0 {
		return ""
	}

	avgMergeTime := totalMergeTime / time.Duration(prCount)
	return fmt.Sprintf("🧮 Total Merged PRs: %d\n⏱ Average Time to Merge: %s", totalPRs, avgMergeTime.Round(time.Minute))
}

func CalculateTotalPRs(token string, since time.Time) (int, error) {
	total := 0
	repos, err := githubclient.GetUserRepos(token, since)
	if err != nil {
		log.Fatalf("Failed to fetch repos: %v", err)
	}

	for _, repo := range repos {
		parts := strings.Split(repo.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner := parts[0]
		name := parts[1]

		prs, err := githubclient.GetUserMergedPRs(token, owner, name, since)
		if err != nil {
			log.Printf("⚠️ Failed to fetch PRs for %s: %v", repo.FullName, err)
			continue
		}

		total += len(prs)
	}

	return total, nil
}

func BuildSummary(client *github.Client, token string, since time.Time) string {
	username, err := githubclient.GetGitHubUsername(token)
	if err != nil {
		log.Println("⚠️ Could not retrieve GitHub username:", err)
		return "No data available"
	}

	repos, err := githubclient.GetPRReposFromSearch(client, username, since)
	if err != nil {
		log.Println("⚠️ Error fetching repos:", err)
		return "No data available"
	}

	totalPRs := 0
	totalCommits := 0
	totalIssues := 0
	repoCount := len(repos)
	stars := 0
	forks := 0

	for _, repo := range repos {
		parts := strings.Split(repo.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repoName := parts[0], parts[1]

		prs, err := githubclient.GetUserMergedPRs(token, owner, repoName, since)
		if err == nil {
			totalPRs += len(prs)
		}

		commitCount, err := githubclient.GetUserCommits(token, owner, repoName, username, since)
		if err == nil {
			totalCommits += commitCount
		}

		totalIssues += repo.IssueFixCount
		stars += repo.StargazersCount
		forks += repo.ForksCount
	}

	return fmt.Sprintf(
		"📦 Repositories: %d\n🟢 PRs Merged: %d\n🔢 Commits: %d\n🐞 Issues Fixed: %d\n⭐ Stars: %d\n🍴 Forks: %d",
		repoCount, totalPRs, totalCommits, totalIssues, stars, forks,
	)
}

func BuildDetailedReport(token string, since time.Time) string {
	username, err := githubclient.GetGitHubUsername(token)
	if err != nil {
		log.Println("⚠️ Could not retrieve GitHub username:", err)
		return "No data available"
	}

	repos, err := githubclient.GetUserRepos(token, since)
	if err != nil {
		log.Println("⚠️ Error fetching repos for detailed report:", err)
		return "No data available"
	}

	var report strings.Builder
	report.WriteString("📊 Developer Metrics Report\n\n")
	report.WriteString("📁 Repositories Included:\n")

	totalCommits := 0
	totalPRs := 0
	var totalMergeTime time.Duration
	prCount := 0

	for _, repo := range repos {
		report.WriteString(fmt.Sprintf(" - %s\n", repo.FullName))

		parts := strings.Split(repo.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repoName := parts[0], parts[1]

		// Languages
		langs, err := githubclient.GetRepoLanguages(token, owner, repoName)
		if err == nil && len(langs) > 0 {
			var langList []string
			for lang := range langs {
				langList = append(langList, lang)
			}
			report.WriteString(fmt.Sprintf("   Languages: %s\n", strings.Join(langList, ", ")))
		}

		// Pull Requests
		prs, err := githubclient.GetUserMergedPRs(token, owner, repoName, since)
		if err == nil {
			for _, pr := range prs {
				report.WriteString(fmt.Sprintf("   🟢 PR: %s\n", pr.Title))
				report.WriteString(fmt.Sprintf("     Description : %s\n", pr.Body))
				report.WriteString(fmt.Sprintf("     📁 Files changed: %d\n", pr.ChangedFiles))
				report.WriteString(fmt.Sprintf("     ✍️ Lines changed: +%d -%d\n", pr.Additions, pr.Deletions))
				totalPRs++
				createdAt, err1 := time.Parse(time.RFC3339, pr.CreatedAt)
				mergedAt, err2 := time.Parse(time.RFC3339, pr.MergedAt)
				if err1 == nil && err2 == nil {
					totalMergeTime += mergedAt.Sub(createdAt)
					prCount++
				}
			}
		}

		// Commits
		commitCount, err := githubclient.GetUserCommits(token, owner, repoName, username, since)
		if err == nil {
			totalCommits += commitCount
		}
	}

	// Pull Request Metrics
	report.WriteString("\n📊 Pull Request Metrics:\n")
	report.WriteString(fmt.Sprintf("🧮 Total Merged PRs: %d\n", totalPRs))
	if prCount > 0 {
		report.WriteString(fmt.Sprintf("⏱ Average Time to Merge: %s\n", (totalMergeTime / time.Duration(prCount)).Round(time.Minute)))
	}

	// Commit-Level Metrics
	report.WriteString("\n📈 Commit-Level Metrics:\n")
	report.WriteString(fmt.Sprintf("🔢 Total Commits: %d\n", totalCommits))

	// Issue Engagement Metrics
	report.WriteString("\n📌 Issue Engagement Metrics:\n")
	report.WriteString(generateIssueEngagementMetrics(repos) + "\n")

	// Collaboration Metrics
	report.WriteString("\n👥 Collaboration Metrics:\n")
	report.WriteString(generateCollaborationMetrics(token, repos, since) + "\n")

	return report.String()
}

func GenerateFullMetricsReport(token string, since time.Time) string {
	repos, err := githubclient.GetUserRepos(token, since)
	if err != nil {
		log.Println("⚠️ Error fetching repos:", err)
		return "No data available"
	}

	var report strings.Builder
	report.WriteString("📊 Full Developer Metrics Report\n\n")
	report.WriteString(fmt.Sprintf("Repositories updated since %s:\n", since.Format("2006-01-02")))
	for _, repo := range repos {
		parts := strings.Split(repo.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repoName := parts[0], parts[1]
		report.WriteString(fmt.Sprintf(" - %s\n", repo.FullName))
		report.WriteString(generateRepoLevelMetrics(token, owner, repoName, since))
		report.WriteString("\n")
	}

	report.WriteString("\n📊 Pull Request Metrics:\n")
	report.WriteString(generatePullRequestMetrics(token, repos, since) + "\n")

	report.WriteString("\n📈 Commit-Level Metrics:\n")
	report.WriteString(generateCommitLevelMetrics(token, repos, since) + "\n")

	report.WriteString("\n📌 Issue Engagement Metrics:\n")
	report.WriteString(generateIssueEngagementMetrics(repos) + "\n")

	report.WriteString("\n👥 Collaboration Metrics:\n")
	report.WriteString(generateCollaborationMetrics(token, repos, since) + "\n")

	return report.String()
}
