package utils

import (
	"os"
	"strings"
)

func AddBadgeToReadme(badgeMarkdown string, readmePath string) error {
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), badgeMarkdown) {
		// Badge already exists
		return nil
	}

	newContent := string(content) + "\n\n" + badgeMarkdown + "\n"
	return os.WriteFile(readmePath, []byte(newContent), 0644)
}
