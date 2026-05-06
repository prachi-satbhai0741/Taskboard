# Project Phases

Each phase is a self-contained Git commit. The history is readable as a learning log.

---

## Phase 1 — Go REST API (Pre-Docker)

**Commit:** `feat: Go REST API with Postgres — pre-Docker baseline`

**Goal:** Get a working API before adding any Docker complexity. Validate the app logic in isolation.

**What's built:**
- Go module initialized
- Gin HTTP router with 5 endpoints
- GORM connected to Postgres
- Task model with soft delete
- `.env` file for local config
- Postgres run via a single `docker run` command (just for the DB, no Compose yet)

**Docker concepts introduced:** None — this is the baseline.

**Test it:**
```bash
go run main.go
curl http://localhost:8080/health
```

---

## Phase 2 — Basic Dockerfile & `.dockerignore`

**Commit:** `feat: basic single-stage Dockerfile and .dockerignore`

**Goal:** Containerize the Go app manually. Understand image layers and build context.

**What's added:**
- Single-stage `Dockerfile`
- `.dockerignore` excluding `.env`, `secrets/`, binaries
- Manual `docker build` and `docker run`

**Docker concepts:** `FROM`, `WORKDIR`, `COPY`, `RUN`, `EXPOSE`, `CMD`, build context, `.dockerignore`

**Commands to run:**
```bash
docker build -t taskflow-go:phase2 .
docker run -p 8080:8080 --env-file .env taskflow-go:phase2
docker images           # observe image size (~300MB+)
docker history taskflow-go:phase2   # inspect layers
```

---

## Phase 3 — Multi-Stage Build & Layer Caching

**Commit:** `feat: multi-stage Dockerfile with layer cache optimization`

**Goal:** Shrink the image significantly. Understand why instruction order matters for caching.

**What's changed:**
- `Dockerfile` refactored to two stages: `builder` and final `alpine`
- `go.mod`/`go.sum` copied before source code for cache efficiency
- `CGO_ENABLED=0` for a fully static binary

**Docker concepts:** multi-stage builds, `COPY --from`, layer caching, `CGO_ENABLED=0`, static binaries

**Commands to run:**
```bash
docker build -t taskflow-go:phase3 .
docker images           # compare size: phase2 vs phase3
# Make a small code change and rebuild — observe which layers are cached
docker build -t taskflow-go:phase3 .
```

---

## Phase 4 — Custom Bridge Network

**Commit:** `feat: connect API and Postgres via custom Docker bridge network`

**Goal:** Run two containers that communicate using service names, not IPs.

**What's added:**
- Custom network `taskflow-net` created manually
- API and Postgres containers both attached to `taskflow-net`
- DB host in `.env` changed from `localhost` to `postgres`

**Docker concepts:** bridge networks, DNS resolution inside user-defined networks, `--network`, `docker network` commands

**Commands to run:**
```bash
docker network create taskflow-net
docker network ls
docker run -d --network taskflow-net --name postgres \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=taskflow postgres:15-alpine
docker run -d --network taskflow-net --name api \
  --env-file .env -p 8080:8080 taskflow-go:phase3
docker network inspect taskflow-net   # see both containers in the network
```

---

## Phase 5 — Docker Compose, Volumes & Health Checks

**Commit:** `feat: Docker Compose with Postgres + Redis, volumes, and health checks`

**Goal:** Replace manual `docker run` commands with a single declarative Compose file. Add data persistence and startup ordering.

**What's added:**
- `docker-compose.yml` with all 3 services
- Named volume `pg-data` for Postgres persistence
- Health check on Postgres so API waits for it to be ready
- Redis service added to the network

**Docker concepts:** Compose services, `depends_on` with `condition: service_healthy`, named volumes, `healthcheck`, multi-service networking

**Commands to run:**
```bash
docker compose up --build
docker compose ps
docker compose logs -f api
docker compose down         # data persists
docker compose down -v      # data deleted
```

---

## Phase 6 — Env Vars, Secrets & Resource Limits

**Commit:** `feat: env_file, Docker secrets, and container resource limits`

**Goal:** Externalize all configuration. Handle sensitive values properly. Constrain container resource usage.

**What's added:**
- `.env.example` committed, `.env` in `.gitignore`
- `env_file` used in Compose instead of hardcoded values
- `secrets/db_password.txt` mounted as a Docker secret
- CPU and memory limits on all containers
- `docker-compose.override.yml` for dev-specific settings (exposed ports, verbose logging)

**Docker concepts:** `env_file`, `secrets`, `deploy.resources.limits`, override files

**Commands to run:**
```bash
docker inspect taskflow-api   # verify env vars are NOT shown in plain text for secrets
docker stats                  # observe CPU/memory limits in real time
```

---

## Phase 7 — Image Tagging & DockerHub Push

**Commit:** `feat: semantic image tagging and DockerHub publish`

**Goal:** Version and publish the image so anyone can pull and run it.

**What's added:**
- Image tagged as `v1.0` and `latest`
- Pushed to DockerHub
- `README.md` updated with pull command

**Docker concepts:** image tagging, semantic versioning, `docker login`, `docker push`, public registries

**Commands to run:**
```bash
docker build -t prachi-satbhai0741/taskflow-go:v1.0 .
docker build -t prachi-satbhai0741/taskflow-go:latest .
docker login
docker push prachi-satbhai0741/taskflow-go:v1.0
docker push prachi-satbhai0741/taskflow-go:latest

# Verify — pull from a clean state
docker rmi prachi-satbhai0741/taskflow-go:latest
docker pull prachi-satbhai0741/taskflow-go:latest
docker run -p 8080:8080 prachi-satbhai0741/taskflow-go:latest
```
