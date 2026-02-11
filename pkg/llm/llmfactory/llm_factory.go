package llmfactory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ionut-t/ask/pkg/llm"
	"github.com/ionut-t/ask/pkg/llm/genai"
	g "google.golang.org/genai"
)

var (
	ErrNoProviderConfigured = errors.New("no LLM provider configured")
	ErrInvalidProvider      = errors.New("unsupported LLM provider")
	ErrMissingCredentials   = errors.New("missing provider credentials")
)

type credentials struct {
	geminiAPIKey    string
	vertexProjectID string
	vertexLocation  string
}

func loadCredentials() credentials {
	return credentials{
		geminiAPIKey:    os.Getenv("GEMINI_API_KEY"),
		vertexProjectID: os.Getenv("VERTEXAI_PROJECT_ID"),
		vertexLocation:  os.Getenv("VERTEXAI_LOCATION"),
	}
}

func New(ctx context.Context, model, provider string, timeout time.Duration) (llm.LLM, error) {
	creds := loadCredentials()

	provider, err := resolveProvider(provider, creds)
	if err != nil {
		return nil, err
	}

	config, err := clientConfig(provider, creds)
	if err != nil {
		return nil, err
	}

	return genai.New(ctx, model, config, timeout)
}

func resolveProvider(provider string, creds credentials) (string, error) {
	if provider != "" {
		return strings.ToLower(strings.TrimSpace(provider)), nil
	}

	switch {
	case creds.geminiAPIKey != "":
		return "gemini", nil
	case creds.vertexProjectID != "" && creds.vertexLocation != "":
		return "vertexai", nil
	default:
		return "", fmt.Errorf("%w: set GEMINI_API_KEY or both VERTEXAI_PROJECT_ID and VERTEXAI_LOCATION", ErrNoProviderConfigured)
	}
}

func clientConfig(provider string, creds credentials) (g.ClientConfig, error) {
	switch provider {
	case "gemini":
		if creds.geminiAPIKey == "" {
			return g.ClientConfig{}, fmt.Errorf("%w for Gemini: GEMINI_API_KEY not set", ErrMissingCredentials)
		}
		return g.ClientConfig{
			Backend: g.BackendGeminiAPI,
			APIKey:  creds.geminiAPIKey,
		}, nil
	case "vertexai":
		var missing []string
		if creds.vertexProjectID == "" {
			missing = append(missing, "VERTEXAI_PROJECT_ID")
		}
		if creds.vertexLocation == "" {
			missing = append(missing, "VERTEXAI_LOCATION")
		}
		if len(missing) > 0 {
			return g.ClientConfig{}, fmt.Errorf("%w for Vertex AI: %s not set", ErrMissingCredentials, strings.Join(missing, " and "))
		}
		return g.ClientConfig{
			Backend:  g.BackendVertexAI,
			Project:  creds.vertexProjectID,
			Location: creds.vertexLocation,
		}, nil
	default:
		return g.ClientConfig{}, fmt.Errorf("%w: %s (supported: gemini, vertexai)", ErrInvalidProvider, provider)
	}
}
