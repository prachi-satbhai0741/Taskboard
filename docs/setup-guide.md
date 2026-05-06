# Setup Guide

## Prerequisites

| Tool | Version | Install |
|---|---|---|
| Go | 1.22+ | [go.dev/dl](https://go.dev/dl/) |
| Docker | 24+ | [docs.docker.com/get-docker](https://docs.docker.com/get-docker/) |
| Docker Compose | v2 (bundled with Docker Desktop) | Included with Docker |
| Git | Any | [git-scm.com](https://git-scm.com/) |

Verify installs:

```bash
go version
docker --version
docker compose version
```

---

## Option 1 — Run Locally (Phase 1, no Docker)

Use this during Phase 1 to run the Go app directly on your machine.

```bash
# Clone
git clone https://github.com/prachi-satbhai0741/taskflow-go.git
cd taskflow-go

# Copy env file
cp .env.example .env
# Open .env and set DB_HOST=localhost

# Start only Postgres via Docker (just for the DB)
docker run -d \
  --name taskflow-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=taskflow \
  -p 5432:5432 \
  postgres:15-alpine

# Install Go dependencies
go mod tidy

# Run the API
go run main.go
```

App runs at `http://localhost:8080`

---

## Option 2 — Run with Docker Compose (Phase 5+)

Use this once `docker-compose.yml` is written in Phase 5.

```bash
# Copy and configure env
cp .env.example .env

# Set up Docker secret
mkdir -p secrets
echo "your_db_password" > secrets/db_password.txt

# Start all services (API + Postgres + Redis)
docker compose up --build

# Or run in background
docker compose up -d --build
```

### Useful Compose Commands

```bash
# View running containers
docker compose ps

# Follow API logs
docker compose logs -f api

# Restart a single service
docker compose restart api

# Stop all containers (keeps volumes)
docker compose down

# Stop and delete all data (wipes DB)
docker compose down -v

# Rebuild image without cache
docker compose build --no-cache
```

---

## Option 3 — Run from DockerHub (Phase 7+)

Once the image is published:

```bash
docker pull prachi-satbhai0741/taskflow-go:latest
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=taskflow \
  prachi-satbhai0741/taskflow-go:latest
```

---

## Environment Variables

Copy `.env.example` to `.env` before running:

```env
DB_HOST=postgres        # use 'localhost' for local run, 'postgres' for Compose
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=            # leave blank if using Docker secrets
DB_NAME=taskflow
PORT=8080
```

**Never commit `.env`** — it is in `.gitignore`. Commit only `.env.example`.

---

## Verify Everything Works

```bash
# Health check
curl http://localhost:8080/health

# Expected response
{"status":"ok","service":"taskflow-go"}
```

If you see this, the API is running and connected to the database correctly.

---

## Troubleshooting

**`connection refused` on port 5432**
Postgres container is not running or not yet healthy. Run `docker compose ps` to check status. Wait for the health check to pass.

**`cannot find module` errors**
Run `go mod tidy` to sync dependencies.

**Port 8080 already in use**
Another process is using the port. Find and stop it: `lsof -i :8080` (Linux/macOS) or change `PORT` in `.env`.

**Compose starts but API exits immediately**
Run `docker compose logs api` to see the error. Usually a missing env var or DB connection failure.
