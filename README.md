
# goploy

A Heroku-lite open-source platform to deploy backend and frontend apps on your VPS with a simple Go CLI.

## GitHub
This project is open-source and published at: [https://github.com/DSwithSiam/goploy](https://github.com/DSwithSiam/goploy)

## Features
- CLI tool with `deploy`, `init`, `destroy`, `logs` commands
- Supports Django, FastAPI, Node, and static frontends
- Auto-generates Nginx and systemd configs
- SSH-based remote deployment


## Installation

1. Install Go (if not already): https://go.dev/doc/install
2. Clone the repo and build the CLI:

```sh
git clone https://github.com/DSwithSiam/goploy.git
cd goploy
go build -o goploy ./cmd/goploy
```

## CLI Commands

- `deploy`   : Deploy an app to your VPS (main workflow)
- `init`     : (coming soon) Initialize a new project
- `destroy`  : (coming soon) Remove a deployed app
- `logs`     : (coming soon) View logs from your app

## How to Deploy with goploy

1. Build the CLI as above.
2. Run the deploy command with your app details:

	**Django Example:**
	```sh
	./goploy deploy --domain example.com --server 1.2.3.4 --project /home/ubuntu/myapp --framework django
	```

	**FastAPI Example:**
	```sh
	./goploy deploy --domain api.example.com --server 1.2.3.4 --project /home/ubuntu/api --framework fastapi
	```

3. Follow the prompts:
	- Enter your SSH username and password for the VPS
	- Enter any environment variables (or leave blank)
	- Answer "y" if you want to deploy a Celery worker

4. The tool will:
	- Connect to your server
	- Install Nginx, Redis, Python dependencies
	- Generate and upload Nginx config
	- Generate and upload systemd service files (Gunicorn/Uvicorn, Celery)
	- Upload your .env file
	- Restart all services

## Example Workflow

```sh
# 1. Build the CLI
git clone https://github.com/DSwithSiam/goploy.git
cd goploy
go build -o goploy ./cmd/goploy

# 2. Deploy your Django app
./goploy deploy --domain example.com --server 1.2.3.4 --project /home/ubuntu/myapp --framework django

# 3. Follow prompts for SSH, .env, and Celery
```

---

## Project Structure
- `cmd/goploy/` - CLI entrypoint
- `internal/cmd/` - Cobra commands
- `pkg/` - Shared utilities

---

For more details, see the code and issues on GitHub.
