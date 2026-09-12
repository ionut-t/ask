package config

import (
	"flag"
	"time"
)

const (
	defaultModel   = "gemini-2.5-flash"
	defaultTimeout = 30 * time.Second
)

type Config struct {
	Model    string
	Provider string
	Timeout  time.Duration
	Copy     bool
}

func ParseConfig() Config {
	var cfg Config

	flag.StringVar(&cfg.Model, "model", defaultModel, "Model to use")
	flag.StringVar(&cfg.Model, "m", defaultModel, "Model to use (shorthand)")
	flag.StringVar(&cfg.Provider, "provider", "", "LLM provider (gemini, vertexai). Auto-detected if not set.")
	flag.DurationVar(&cfg.Timeout, "timeout", defaultTimeout, "Request timeout")
	flag.DurationVar(&cfg.Timeout, "t", defaultTimeout, "Request timeout (shorthand)")
	flag.BoolVar(&cfg.Copy, "copy", false, "Copy response to clipboard")
	flag.BoolVar(&cfg.Copy, "c", false, "Copy response to clipboard (shorthand)")

	flag.Parse()

	return cfg
}
