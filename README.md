# goploy

A Heroku-lite open-source platform to deploy backend and frontend apps on your VPS with a simple Go CLI.

## Features
- CLI tool with `deploy`, `init`, `destroy`, `logs` commands
- Supports Django, FastAPI, Node, and static frontends
- Auto-generates Nginx and systemd configs
- SSH-based remote deployment

## Usage
```
goploy deploy --domain example.com --server 1.2.3.4 --project /path/to/app --framework django
```

## Project Structure
- `cmd/goploy/` - CLI entrypoint
- `internal/cmd/` - Cobra commands
- `pkg/` - Shared utilities
