package service

import (
	"fmt"
	"log"
	"os"
	"pm/client"
	"strings"
	"time"
)

func GetBadgesFromReadme(readmePath string) ([]string, []string, error) {

	content, err := getReadMeContent(readmePath)
	if err != nil {
		log.Fatalf("Failed to read README: %v", err)
	}

	existingBadges := []string{}
	newBadges := []string{}

	if strings.Contains(content, "1st%20PR-achieved") {
		existingBadges = append(existingBadges, "🎉 1st PR Badge")
	}
	if strings.Contains(content, "1st%20Repo-active") {
		existingBadges = append(existingBadges, "📁 1st Repo Badge")
	}

	return existingBadges, newBadges, nil
}

func getReadMeContent(readmePath string) (string, error) {
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		return "", fmt.Errorf("failed to read README: %w", err)
	}
	return string(readmeContent), nil
}

func DisplayExistingBadges(badges []string, label string) {
	fmt.Printf("\n%s\n", label)
	if len(badges) == 0 {
		fmt.Println("- None yet")
	} else {
		for _, badge := range badges {
			fmt.Println("-", badge)
		}
	}
}

func firstPrBadge(token string, since time.Time, readmePath string) bool {
	// Calculate total PRs in the time window
	totalPRs, err := CalculateTotalPRs(token, since)
	if err != nil {
		log.Printf("Error calculating total PRs: %v", err)
		return false
	}

	// Read README content
	readMeContent, err := getReadMeContent(readmePath)
	if err != nil {
		log.Printf("Failed to read README: %v", err)
		// Still allow badge if README can't be read — up to you
		return totalPRs >= 1 // optional: you could choose to return false here instead
	}

	// Badge condition: user has at least 1 PR and doesn't already have badge
	return totalPRs >= 1 && !strings.Contains(readMeContent, "1st%20PR-achieved")
}

func firstRepoBadge(token string, since time.Time, readmePath string) bool {
	userRepos, err := client.GetUserRepos(token, since)
	if err != nil {
		log.Printf("Error getting user repos: %v", err)
		return false
	}

	readMeContent, err := getReadMeContent(readmePath)
	if err != nil {
		log.Printf("Failed to read README: %v", err)
		// You can decide to return false here or assume the badge isn't present
		return len(userRepos) >= 1
	}

	return len(userRepos) >= 1 && !strings.Contains(readMeContent, "1st%20Repo-active")
}

func GetNewBadges(token string, badges []string, since time.Time, readmePath string) (newBadges []string) {
	// Determine which badges should be added
	if firstPrBadge(token, since, readmePath) {
		badges = append(badges, "![First PR](https://img.shields.io/badge/🎉%201st%20PR-achieved-green)")
	}

	if firstRepoBadge(token, since, readmePath) {
		badges = append(badges, "![First Repo](https://img.shields.io/badge/📁%201st%20Repo-active-blue)")

	}
	return badges
}
