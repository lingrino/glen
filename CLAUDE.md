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
2. `Variables.Init()` creates a GitLab API client and fetches variables, respecting [GitLab's variable precedence](https://docs.gitlab.com/ee/ci/variables/#priority-of-environment-variables) (group vars fetched first, then project vars override them)
3. Output is formatted as `export` (shell-ready), `json`, or `table`

**Key CLI flags:**
- `-k, --api-key` - GitLab API key (defaults to `GITLAB_TOKEN` env var)
- `-r, --recurse` - Include variables from parent groups (defaults to false, so only project vars by default)
- `-g, --group-only` - Only fetch group variables, skip project variables
- `-o, --output` - Output format: `export` (default), `json`, or `table`
- `-d, --directory` - Git repo path (defaults to current directory)
- `-n, --remote-name` - Git remote name (defaults to `origin`)

**Key behaviors:**
- Supports both SSH (`git@`) and HTTP(S) remote URLs
- Parent groups are processed from outermost to innermost, with project variables having highest precedence per GitLab's priority rules
