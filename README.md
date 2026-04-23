
# goploy

A Heroku-lite open-source platform to deploy backend and frontend apps on your VPS with a simple Go CLI.

## GitHub
This project is open-source and published at: [https://github.com/DSwithSiam/goploy](https://github.com/DSwithSiam/goploy)

## Features
- CLI tool with `deploy`, `init`, `destroy`, `logs` commands
- Supports Django, FastAPI, Node, and static frontends
- Auto-generates Nginx and systemd configs
- SSH-based remote deployment



## How to Install & Use Goploy

### 1. Traditional way (clone & build)

```sh
git clone https://github.com/DSwithSiam/goploy.git
cd goploy
go build -o goploy ./cmd/goploy
./goploy --help
```

**Problems:**
- Requires Go and build tools
- More dependencies, more errors possible
- Not as user-friendly for non-developers

---

### 2. Release-based install (easy way)

**Recommended for most users!**

**Step 1: Download the latest release**

Go to: https://github.com/DSwithSiam/goploy/releases
Download `goploy-linux-amd64.tar.gz` (or the binary for your OS)

Or use wget:
```sh
wget https://github.com/DSwithSiam/goploy/releases/download/v1.0.0/goploy-linux-amd64
```

**Step 2: Extract and install**
```sh
tar -xvzf goploy-linux-amd64.tar.gz
cd goploy
chmod +x goploy
sudo mv goploy /usr/local/bin/
```

**Step 3: Run**
```sh
goploy --help
```

---

## Usage Example

Deploy a Django app:
```sh
goploy deploy --domain example.com --server 1.2.3.4 --project /home/ubuntu/myapp --framework django
```

Deploy a FastAPI app:
```sh
goploy deploy --domain api.example.com --server 1.2.3.4 --project /home/ubuntu/api --framework fastapi
```

You’ll be prompted for SSH credentials, .env variables, and Celery worker (optional).


### Deploy a Static Frontend (HTML/JS/CSS)

Suppose your static site is in `/home/ubuntu/frontend` and you want to serve it at `frontend.example.com`:

```sh
goploy deploy --domain frontend.example.com --server 1.2.3.4 --project /home/ubuntu/frontend --framework static
```

**What happens:**
- Goploy will set up Nginx to serve your static files from `/home/ubuntu/frontend`.
- No backend server (Gunicorn/Uvicorn) is needed for static sites.

**You’ll be prompted for:**
- SSH credentials
- (You can skip .env and Celery prompts for static sites)


---

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



## Contributing
PRs and issues welcome! See CONTRIBUTING.md.

---

## License
MIT
