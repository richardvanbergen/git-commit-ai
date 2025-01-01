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

func (a AnthropicAdapter) Summerize(query string) string {
	client := anthropic.NewClient(a.Key)
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model:     anthropic.ModelClaude3Dot5Sonnet20241022,
		MaxTokens: 1000,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(query),
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
