package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CommitAndPushProfileReadme() {
	profileRepoPath := "/Users/yaswood/yassir20191"

	commands := [][]string{
		{"git", "-C", profileRepoPath, "add", "README.md"},
		{"git", "-C", profileRepoPath, "commit", "-m", "🤖 Update badges in README"},
		{"git", "-C", profileRepoPath, "push", "origin", "main"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("🔧 Running:", strings.Join(cmdArgs, " "))
		if err := cmd.Run(); err != nil {
			if cmdArgs[1] == "-C" && cmdArgs[3] == "commit" {
				fmt.Println("ℹ️ No changes to commit.")
			} else {
				fmt.Println("❌ Error running command:", err)

			}
		}
	}
}
