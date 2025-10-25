# Claude Code Configuration for linctl

## Project Overview

**linctl** is a comprehensive CLI tool for Linear's API, designed for both agents and humans. This is a fork of [dorkitude/linctl](https://github.com/dorkitude/linctl) maintained by Mike Beaumier with additional features and enhancements.

**Repository**: https://github.com/delphos-mike/linctl
**Language**: Go 1.23+
**Current Version**: v0.2.0

## Key Features

- Issue management (create, list, update, assign, archive)
- Issue relations (blocks, blocked-by, related, duplicate)
- Full-text issue search
- Project tracking and management
- Team and user management
- Comments and attachments
- Webhook configuration
- Multiple output formats (table, plaintext, JSON)

## Project Structure

```
linctl/
├── cmd/                    # Cobra commands (auth, issue, project, team, user, comment)
│   ├── root.go            # Root command and global flags
│   ├── issue.go           # Issue commands (list, create, update, relate, search)
│   ├── project.go         # Project commands
│   ├── team.go            # Team commands
│   ├── user.go            # User commands
│   ├── comment.go         # Comment commands
│   └── docs.go            # Documentation command
├── pkg/                   # Shared packages
│   ├── api/              # Linear API client and queries
│   ├── auth/             # Authentication management
│   ├── output/           # Output formatting (table, JSON, plaintext)
│   └── utils/            # Utility functions (time parsing)
├── docs/                  # Documentation snapshots
├── Formula/               # Homebrew formula (if maintaining tap)
├── main.go               # Entry point
├── Makefile              # Build and development tasks
├── go.mod                # Go module definition
├── README.md             # User documentation
├── CONTRIBUTING.md       # Development and release workflow
└── smoke_test.sh         # Automated smoke tests

```

## Development Workflow

### Building and Testing

```bash
# Install dependencies
make deps

# Build the binary
make build

# Run without building
go run main.go [command]

# Run smoke tests (requires Linear API token)
make test

# Format code
make fmt

# Lint (requires golangci-lint)
make lint
```

### Key Make Targets

- `make build` - Build binary with version injection
- `make clean` - Remove build artifacts
- `make install` - Install to /usr/local/bin (requires sudo)
- `make test` - Run smoke tests (read-only operations)
- `make deps` - Install/update dependencies
- `make fmt` - Format Go code
- `make lint` - Run linter

### Running Locally

```bash
# Set up authentication first
go run main.go auth

# Test commands
go run main.go issue list
go run main.go project list
go run main.go whoami
```

## Architecture

### Command Structure

linctl uses [Cobra](https://github.com/spf13/cobra) for CLI commands and [Viper](https://github.com/spf13/viper) for configuration management.

- **Root command**: `cmd/root.go` - defines global flags (`--json`, `--plaintext`)
- **Subcommands**: Each in `cmd/<resource>.go` (e.g., `cmd/issue.go`)
- **API Client**: `pkg/api/client.go` - GraphQL client for Linear API
- **Queries**: `pkg/api/queries.go` - GraphQL query definitions
- **Auth**: `pkg/auth/auth.go` - API key management (stored in `~/.config/linctl/config.yaml`)

### Adding New Commands

1. Add command definition in appropriate `cmd/<resource>.go` file
2. Define GraphQL query in `pkg/api/queries.go` if needed
3. Implement command logic using `api.Client`
4. Use `output` package for formatted output
5. Add to smoke tests in `smoke_test.sh`
6. Update README with examples

### Output Formats

The tool supports three output formats via global flags:
- **Table** (default): Human-readable tables via `tablewriter`
- **Plaintext**: Simple text output with `--plaintext`
- **JSON**: Machine-readable with `--json`

Use the `output` package functions for consistent formatting.

## Code Style

- Follow standard Go conventions (gofmt)
- Use meaningful variable names
- Add comments for complex logic
- Keep functions focused and single-purpose
- Error messages should be clear and actionable

## Testing

### Smoke Tests

The `smoke_test.sh` script runs read-only commands to verify basic functionality:
- Authentication check
- Issue listing
- Project listing
- User info
- Team listing

Run with: `make test` or `./smoke_test.sh`

### Manual Testing

When adding features:
1. Test with various flag combinations
2. Verify JSON output is valid
3. Check error handling for invalid inputs
4. Test against your Linear workspace

## Release Process

See `CONTRIBUTING.md` for complete release checklist. Summary:

1. **Prepare**: Update docs, run tests
2. **Tag**: Create semantic version tag (e.g., `v0.3.0`)
3. **Release**: Push tag and create GitHub release
4. **Update tap**: If maintaining Homebrew tap (optional for this fork)

```bash
# Tag and push
git tag v0.X.Y -a -m "vX.Y.Z: Brief description"
git push origin v0.X.Y

# Create GitHub release
gh release create v0.X.Y \
  --title "linctl v0.X.Y" \
  --notes "Release notes here"
```

Version is automatically injected at build time via Makefile LDFLAGS.

## Fork-Specific Information

### Changes from Upstream

- Added issue relation management (`relate`/`unrelate` commands)
- Enhanced issue search functionality
- Updated module path to `github.com/delphos-mike/linctl`
- Migrated from `master` to `main` branch
- Updated copyright to include Mike Beaumier

### Syncing with Upstream

This fork does not track upstream changes. The `upstream` remote has been removed. If you want to pull updates from the original repo:

```bash
git remote add upstream https://github.com/dorkitude/linctl.git
git fetch upstream
git merge upstream/main  # May require conflict resolution
```

## Authentication

linctl uses Linear Personal API Keys stored in `~/.config/linctl/config.yaml`:

```yaml
linear_api_key: lin_api_xxxxxxxxxxxxxxxxxxxxx
```

Get your API key from: https://linear.app/settings/api

## Common Tasks

### Adding a New Issue Command

1. Edit `cmd/issue.go`
2. Define command using Cobra pattern
3. Add to `init()` function: `issueCmd.AddCommand(newCmd)`
4. Implement GraphQL query in `pkg/api/queries.go`
5. Update smoke tests if applicable

### Modifying API Queries

All GraphQL queries are in `pkg/api/queries.go`. The Linear API schema reference is at:
https://developers.linear.app/docs/graphql/schema

### Debugging

- Use `--json` flag to see raw API responses
- Check `~/.config/linctl/config.yaml` for auth issues
- Set `LINCTL_DEBUG=1` environment variable for verbose logging (if implemented)

## Dependencies

Key dependencies (see `go.mod`):
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/olekukonko/tablewriter` - Table formatting
- `github.com/fatih/color` - Colored output

## Installation Methods

For end users (documented in README.md):
- From source: `make install`
- Via asdf: Using custom plugin (in development)
- Direct Go install: `go install github.com/delphos-mike/linctl@latest`

## Troubleshooting

### Build Issues
- Ensure Go 1.23+ is installed
- Run `go mod download` to fetch dependencies
- Check for import path issues after module updates

### API Issues
- Verify API key in `~/.config/linctl/config.yaml`
- Check Linear API rate limits (5,000 requests/hour)
- Use `--json` to see full error responses

### Version Issues
- Version is set via git tags and Makefile LDFLAGS
- `go run` always shows "dev" (version only injected in built binary)
- Check tags with `git tag -l`

## Links

- GitHub Repository: https://github.com/delphos-mike/linctl
- Original Repository: https://github.com/dorkitude/linctl
- Linear API Docs: https://developers.linear.app/
- Issue Tracker: https://github.com/delphos-mike/linctl/issues
