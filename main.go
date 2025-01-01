package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/richardvanbergen/git-commit-ai/adapters"
)

func getEditor() (string, error) {
	// Try git config first
	cmd := exec.Command("git", "config", "--get", "core.editor")
	if output, err := cmd.Output(); err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output)), nil
	}

	// Fallback to environment variables
	if editor := os.Getenv("GIT_EDITOR"); editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	// Final fallback
	return "vim", nil
}

func editMessage(initialContent string) (string, error) {
	editor, err := getEditor()
	if err != nil {
		return "", fmt.Errorf("failed to get editor: %w", err)
	}

	// Create temporary file
	tmpfile, err := os.CreateTemp("", "COMMIT_EDITMSG.*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write initial content
	if _, err := tmpfile.WriteString(initialContent); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpfile.Close()

	// Open editor
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run editor: %w", err)
	}

	// Read edited content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read edited file: %w", err)
	}

	return string(content), nil
}

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

func commitChanges(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

func confirmCommit(message string) bool {
	fmt.Println("\nProposed commit message:")
	fmt.Println("------------------------")
	fmt.Println(message)
	fmt.Println("------------------------")
	fmt.Print("Do you want to commit with this message? (y/N): ")

	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
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

	if len(stagedFiles) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No files are staged for commit")
		os.Exit(1)
	}

	client := adapters.NewAnthropicClient()

	suggestedMessage := client.Summerize(
		fmt.Sprintf("Changed files:\n%s\n\nDiff:\n%s", strings.Join(stagedFiles, "\n"), diff),
	)

	finalMessage, err := editMessage(suggestedMessage)
	if err != nil {
		fmt.Printf("Error editing message: %v\n", err)
		os.Exit(1)
	}

	if !confirmCommit(finalMessage) {
		fmt.Println("Commit cancelled")
		os.Exit(0)
	}

	if err := commitChanges(finalMessage); err != nil {
		fmt.Printf("Error creating commit: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Changes committed successfully!")
}
