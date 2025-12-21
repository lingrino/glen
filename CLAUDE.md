# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build
go build

# Run tests with coverage
go test -cover -coverprofile=c.out -covermode=atomic -race -v ./...

# Run a single test
go test -v -run TestNewRepo ./glen/

# Lint (uses golangci-lint with extensive linter configuration)
golangci-lint run --timeout=5m

# Check for vulnerabilities
govulncheck ./...

# Verify go.mod is tidy
go mod tidy
```

## Architecture

Glen is a CLI tool that fetches GitLab CI/CD environment variables from a project and its parent groups.

**Package structure:**
- `main.go` - Entry point, passes version (set by goreleaser at build time) to cmd
- `cmd/` - CLI commands using cobra (root command + version subcommand)
- `glen/` - Core library with two main types:
  - `Repo` - Parses local git repo to extract GitLab project path, base URL, and parent group hierarchy
  - `Variables` - Uses GitLab API to fetch CI/CD variables from project and groups

**Data flow:**
1. `Repo.Init()` reads the git remote URL and parses it into GitLab-compatible paths
2. `Variables.Init()` creates a GitLab API client and fetches variables, respecting GitLab's variable precedence (group vars first, then project vars override)
3. Output is formatted as `export` (shell-ready), `json`, or `table`

**Key behaviors:**
- API key comes from `GITLAB_TOKEN` env var by default, or `--api-key` flag
- Supports both SSH (`git@`) and HTTP(S) remote URLs
- Parent groups are processed from outermost to innermost, with project variables having highest precedence
