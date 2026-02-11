package genai

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/genai"
)

type GenAI struct {
	model   string
	client  *genai.Client
	timeout time.Duration
}

func New(ctx context.Context, model string, config genai.ClientConfig, timeout time.Duration) (*GenAI, error) {
	client, err := genai.NewClient(ctx, &config)
	if err != nil {
		return nil, err
	}

	return &GenAI{
		client:  client,
		model:   model,
		timeout: timeout,
	}, nil
}

func (g *GenAI) Generate(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	result, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	if result == nil || len(result.Candidates) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return result.Text(), nil
}

func (g *GenAI) Close() error {
	g.client = nil
	return nil
}
