package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/liushuangls/go-anthropic/v2"
)

type AnthropicAdapter struct {
	Key string
}

func (a AnthropicAdapter) Summerize(changes string) string {
	assitMessage := `
    You are a Git Commit Message Assistant, an AI trained to generate clear, concise, and informative commit messages. Your task is to analyze the provided code changes and context, then create a commit message that follows best practices. Adhere to these guidelines:

    <commit_structure>
    1. Subject line:
      - Start with an imperative verb (e.g., Add, Fix, Update)
      - Limit to 50 characters
      - Capitalize only the first letter
      - No period at the end

    2. Body (if needed):
      - Separate from subject with a blank line
      - Wrap at 72 characters per line
      - Explain the what and why of changes, not how
      - Use bullet points for multiple items
    </commit_structure>

    <additional_instructions>
    - Be specific and descriptive about the changes
    - Reference relevant issue numbers or pull requests
    - Use the imperative mood consistently
    - Avoid redundant or unnecessary information
    - If multiple significant changes, use a bulleted list in the body
    </additional_instructions>

    Analyze the provided code changes and context, then generate a commit message following these guidelines. Your response should only contain the commit message, without any additional explanation or commentary.
  `

	client := anthropic.NewClient(a.Key)
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model:     anthropic.ModelClaude3Dot5Sonnet20241022,
		MaxTokens: 1000,
		Messages: []anthropic.Message{
			anthropic.NewAssistantTextMessage(assitMessage),
			anthropic.NewUserTextMessage(changes),
		},
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return ""
	}

	return resp.Content[0].GetText()
}

func NewAnthropicClient() QueryableLLM {
	apiKey := os.Getenv("CLAUDE_API_KEY")

	if apiKey == "" {
		fmt.Println("CLAUDE_API_KEY environment variable is not set")
		os.Exit(1)
	}

	return AnthropicAdapter{
		Key: apiKey,
	}
}
