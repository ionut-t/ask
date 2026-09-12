package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ionut-t/ask/internal/config"
	"github.com/ionut-t/ask/pkg/llm/llmfactory"
)

const maxInputSize = 1 << 20 // 1 MB

func Execute() {
	cfg := config.ParseConfig()

	prompt, err := getPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(prompt) == "" {
		fmt.Fprintln(os.Stderr, "No input provided. Provide a prompt via stdin or command line arguments.")
		fmt.Fprintln(os.Stderr, "Usage: echo 'Your prompt' | ask --model gemini-2.5-flash")
		os.Exit(1)
	}

	ctx := context.Background()

	llm, err := llmfactory.New(ctx, cfg.Model, cfg.Provider, cfg.Timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initialising LLM: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := llm.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing LLM: %v\n", err)
		}
	}()

	response, err := llm.Generate(ctx, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating response: %v\n", err)
		os.Exit(1)
	}

	if cfg.Copy {
		copyToClipboard(response)
	}

	fmt.Println("\n" + response)
}

func getPrompt() (string, error) {
	var input []string

	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat stdin: %w", err)
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := io.ReadAll(io.LimitReader(os.Stdin, maxInputSize))
		if err != nil {
			return "", err
		}
		if len(bytes) > 0 {
			input = append(input, string(bytes))
		}
	}

	if len(flag.Args()) > 0 {
		input = append(input, strings.Join(flag.Args(), " "))
	}

	if len(input) == 0 {
		return "", nil
	}

	return strings.Join(input, "\n"), nil
}

func copyToClipboard(text string) {
	text = strings.TrimSpace(text)
	// trim backticks if response is wrapped in a code block
	if strings.HasPrefix(text, "```") {
		if end := strings.Index(text, "\n"); end != -1 {
			text = text[end+1:]
		}
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}

	if err := clipboard.WriteAll(text); err != nil {
		fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\n", err)
	}
}
