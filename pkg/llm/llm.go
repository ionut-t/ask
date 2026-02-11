package llm

import "context"

type LLM interface {
	Generate(ctx context.Context, prompt string) (string, error)
	Close() error
}
