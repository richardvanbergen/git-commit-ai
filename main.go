package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/richardvanbergen/git-commit-ai/adapters"
)

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	files := []string{}
	for _, file := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if file != "" {
			files = append(files, file)
		}
	}

	return files, nil
}

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	return string(output), nil
}

func main() {
	stagedFiles, err := GetStagedFiles()
	if err != nil {
		fmt.Printf("Error getting staged files: %v\n", err)
		os.Exit(1)
	}

	diff, err := GetStagedDiff()
	if err != nil {
		fmt.Printf("Error getting staged diff: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Staged files:")
	for _, file := range stagedFiles {
		fmt.Printf("- %s\n", file)
	}

	client := adapters.NewAnthropicClient()

	res := client.Summerize(
		fmt.Sprintf("Changed files:\n%s\n\nDiff:\n%s", strings.Join(stagedFiles, "\n"), diff),
		"Write a great git commit message to summerizes the diffs and file changes. The title should be short but you can include more information in the description if it's necessary.",
	)

	fmt.Println(res)
}
