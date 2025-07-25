package main

import (
	"fmt"
	"log"
	"os"
	gitClient "pm/client"
	gitService "pm/service"
	"pm/tui"
	"pm/utils"
	"strings"
	"time"
)

func main() {
	activeSessionClient := gitClient.NewGitHubClient("GITHUB_TOKEN")
	if len(os.Args) > 1 && os.Args[1] == "tui" {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			log.Fatal("Set GITHUB_TOKEN environment variable.")
		}

		defaultSince := time.Now().AddDate(0, 0, -7)
		defaultSummary := gitService.BuildSummary(activeSessionClient, token, defaultSince)

		period, since, ok := tui.RunWithTokenWithSummary(token, defaultSummary)

		if !ok {
			log.Println("No report generated.")
			return
		}
		if period == "badges" {
			readmePath := "/Users/yaswood/yassir20191/README.md"
			content, err := os.ReadFile(readmePath)
			if err != nil {
				log.Printf("⚠️ Could not read README: %v", err)
			} else {
				existingBadges := []string{}
				text := string(content)

				if strings.Contains(text, "1st%20PR-achieved") {
					existingBadges = append(existingBadges, "🎉 1st PR Badge")
				}
				if strings.Contains(text, "1st%20Repo-active") {
					existingBadges = append(existingBadges, "📁 1st Repo Badge")
				}

				fmt.Println("\n🏅 Developer Badges:")
				if len(existingBadges) == 0 {
					fmt.Println("- No badges found.")
				} else {
					for _, badge := range existingBadges {
						fmt.Println("-", badge)
					}
				}
			}
			return
		}

		summary := gitService.BuildDetailedReport(token, since)
		dateStr := time.Now().Format("2006-01-02")
		os.MkdirAll("reports", os.ModePerm)
		filename := fmt.Sprintf("reports/report_%s.txt", dateStr)
		content := fmt.Sprintf("%s report\n\n%s", strings.Title(period), summary)

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to write report: %v", err)
		}

		fmt.Printf("✅ Report saved to %s\n", filename)
		return
	}

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

	repos, err := gitClient.GetUserRepos(token, sinceDate)
	if err != nil {
		log.Fatal(err)
	}

	report := gitService.GenerateFullMetricsReport(token, sinceDate)
	fmt.Println(report)

	// Badge generation: check if user has at least 1 merged PR
	totalPRs := 0
	for _, repo := range repos {
		parts := strings.Split(repo.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repoName := parts[0], parts[1]
		prs, err := gitClient.GetUserMergedPRs(token, owner, repoName, sinceDate)
		if err != nil {
			continue
		}
		totalPRs += len(prs)
	}

	// --- Begin badge block ---
	readmePath := "/Users/yaswood/yassir20191/README.md"

	existingBadges, newBadges, err := gitService.GetBadgesFromReadme(readmePath)
	gitService.DisplayExistingBadges(existingBadges, "\n📛 Existing Badges in README:")
	newBadges = gitService.GetNewBadges(token, newBadges, sinceDate, readmePath)

	fmt.Println("\n🏅 Newly Unlocked Badges:")
	if len(newBadges) == 0 {
		fmt.Println("No new badges to add.")
	} else {
		for _, badge := range newBadges {
			fmt.Println("- Adding:", badge)
			err := utils.AddBadgeToReadme(badge, readmePath)
			if err != nil {
				log.Println("❌ Failed to update README with badge:", err)
			} else {
				fmt.Println("✅ Badge added to README.")
			}
		}
	}
	// --- End badge block ---

	if len(newBadges) > 0 {
		utils.CommitAndPushProfileReadme()
	}
}
