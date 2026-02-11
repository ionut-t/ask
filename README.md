# ask

A CLI tool to query Google LLMs (Gemini API / Vertex AI) from the terminal.

## Installation

```bash
go install github.com/ionut-t/ask@latest
```

## Configuration

Set environment variables:

**Gemini API:**

```bash
export GEMINI_API_KEY="api-key"
```

**Vertex AI:**

```bash
export VERTEXAI_PROJECT_ID="project-id"
export VERTEXAI_LOCATION="us-central1"
```

If both are set, Gemini is used by default. Use `--provider` to override.

## Usage

```bash
echo "What is Go?" | ask
ask "What is Go?"
cat main.go | ask -m gemini-2.5-pro "Find bugs in this code"
git diff HEAD~1 | ask "Write a commit message for this diff"
```

## Flags

| Flag         | Short | Default            | Description                         |
| ------------ | ----- | ------------------ | ----------------------------------- |
| `--model`    | `-m`  | `gemini-2.5-flash` | Model to use                        |
| `--provider` |       | auto-detected      | LLM provider (`gemini`, `vertexai`) |
| `--timeout`  | `-t`  | `30s`              | Request timeout                     |
