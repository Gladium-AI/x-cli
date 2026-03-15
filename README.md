# x-cli

A command-line tool for interacting with X (Twitter) via its internal GraphQL API.

Fetch timelines, tweets, user profiles, search results, followers, and following lists directly from your terminal.

## Installation

**Prerequisites:** Go 1.25+ and Google Chrome (for browser-based login).

```bash
git clone https://github.com/paolo/x-cli.git
cd x-cli
make build
sudo make install
```

This installs `x-cli` to `/usr/local/bin`. To uninstall:

```bash
sudo make uninstall
```

## Authentication

x-cli uses browser-based authentication. On login, a Chrome window opens for you to sign in to X (supports Google OAuth, email/password, etc.). Session credentials are stored locally at `~/.x-cli/credentials.json`.

```bash
# Log in via browser
x-cli auth login

# Check auth status
x-cli auth status

# Log out (clears stored credentials)
x-cli auth logout
```

## Commands

### Timeline

```bash
# Home timeline
x-cli timeline home

# A user's tweets
x-cli timeline user @elonmusk

# Paginate through all results
x-cli timeline home --all --max-pages 5
```

| Flag | Default | Description |
|------|---------|-------------|
| `--count` | 20 | Tweets per page |
| `--cursor` | | Pagination cursor |
| `--all` | false | Auto-paginate through results |
| `--max-pages` | 10 | Max pages when using `--all` |

### Search

```bash
# Top results
x-cli search "golang"

# Latest tweets
x-cli search "breaking news" --type latest

# People
x-cli search "elon" --type people
```

| Flag | Default | Description |
|------|---------|-------------|
| `--count` | 20 | Results per page |
| `--cursor` | | Pagination cursor |
| `--type` | top | Search type: `top`, `latest`, `people`, `media` |

### User Profile

```bash
x-cli user get @paoloanzn
```

### Tweet

Accepts a tweet ID or full URL:

```bash
x-cli tweet get 1234567890
x-cli tweet get https://x.com/user/status/1234567890
```

### Followers & Following

```bash
x-cli followers @paoloanzn
x-cli following @paoloanzn --count 50
```

| Flag | Default | Description |
|------|---------|-------------|
| `--count` | 20 | Users per page |
| `--cursor` | | Pagination cursor |

## Global Flags

These flags work with every command:

| Flag | Description |
|------|-------------|
| `--json` | Output raw JSON from the API |
| `--verbose` | Print request URLs, HTTP status codes, and response details |

## Pagination

Most commands that return lists support cursor-based pagination. When results are available, x-cli prints a "Next page" hint with the command to run:

```
Next page: x-cli timeline home --cursor "DAABCgABGRI..."
```

Use `--all` on timeline commands to auto-paginate. Rate limits are respected automatically with wait-and-retry.

## Project Structure

```
x-cli/
├── main.go
├── Makefile
├── cmd/                    # CLI command definitions (cobra)
│   ├── root.go
│   ├── auth.go
│   ├── timeline.go
│   ├── user.go
│   ├── tweet.go
│   ├── search.go
│   └── followers.go
└── internal/
    ├── api/                # GraphQL API client & endpoint registry
    ├── auth/               # Browser login flow & credential storage
    ├── models/             # Domain types (Tweet, User, Timeline)
    └── output/             # Pretty-print & JSON output formatting
```

## Query ID Rotation

X rotates GraphQL query IDs on every deploy. If you start seeing 404 errors across all commands, the IDs need updating. They can be extracted from X's production JavaScript bundle (`main.*.js`) or captured from browser network traffic. Update `internal/api/endpoints.go` with the new IDs.

## License

MIT
