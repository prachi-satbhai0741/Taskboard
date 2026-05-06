# Docker Concepts

Every Docker concept used in this project, explained in context — not in isolation.

---

## Phase 2 — Dockerfile & `.dockerignore`

### Dockerfile

A `Dockerfile` is a text file with instructions Docker reads top-to-bottom to build an image. Each instruction creates a new **layer**.

```dockerfile
FROM golang:1.22-alpine        # base image
WORKDIR /app                   # set working directory inside container
COPY go.mod go.sum ./          # copy dependency files first
RUN go mod download            # download dependencies
COPY . .                       # copy source code
RUN go build -o taskflow .     # compile the binary
EXPOSE 8080                    # document the port (does not publish it)
CMD ["./taskflow"]             # command to run when container starts
```

**Key mental model:** an image is a stack of read-only layers. A container is that stack with a thin writable layer on top. When you `docker build`, Docker checks if each layer is already cached before re-running it.

### `.dockerignore`

Tells Docker what to exclude from the build context sent to the daemon. Without it, your entire project folder (including `node_modules`, `.env`, secrets) gets sent.

```
.env
secrets/
*.md
.git/
tmp/
```

---

## Phase 3 — Multi-Stage Build & Layer Caching

### Multi-Stage Build

Uses multiple `FROM` instructions in one Dockerfile. The final image only contains what the last stage copies — not the compiler, not the source code.

```dockerfile
# Stage 1 — Builder
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o taskflow .

# Stage 2 — Final image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/taskflow .   # copies only the binary
EXPOSE 8080
CMD ["./taskflow"]
```

**Result:** final image is ~20MB instead of ~300MB. The Go compiler never ships to production.

### Layer Caching

Docker caches each layer. If a layer's instruction and its inputs haven't changed, Docker reuses the cached layer instead of re-running it.

**Bad order (breaks cache on every code change):**
```dockerfile
COPY . .             # copies everything — invalidates cache on any file change
RUN go mod download  # re-runs every time even if go.mod didn't change
```

**Good order (cache-friendly):**
```dockerfile
COPY go.mod go.sum ./   # only changes when dependencies change
RUN go mod download     # cached unless go.mod changes
COPY . .                # source code copied after — only this layer rebuilds
RUN go build ...
```

---

## Phase 4 — Custom Bridge Networks

By default, containers can't talk to each other unless they're on the same network. Docker's default bridge network (`docker0`) does not support DNS-based container name resolution.

A **user-defined bridge network** does — containers reference each other by service name.

```bash
# Create network
docker network create taskflow-net

# Run containers on the same network
docker run --network taskflow-net --name postgres postgres:15-alpine
docker run --network taskflow-net --name api taskflow-go
```

Inside the API container, the DB host is simply `postgres` — Docker resolves it to the correct container IP automatically.

**Why not use `--link`?** It's deprecated. User-defined networks are the correct approach.

---

## Phase 5 — Volumes, Docker Compose & Health Checks

### Volumes

Containers are ephemeral — their filesystem is destroyed when they stop. A **volume** persists data outside the container lifecycle.

```yaml
volumes:
  pg-data:              # named volume managed by Docker

services:
  postgres:
    volumes:
      - pg-data:/var/lib/postgresql/data   # mount volume into container
```

`pg-data` survives `docker compose down` but is removed with `docker compose down -v`.

### Docker Compose

A YAML file that declares all services, networks, volumes, and config in one place — so you don't run 10 `docker run` commands manually.

```yaml
services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: taskflow
    volumes:
      - pg-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pg-data:
```

### Health Checks

`depends_on` alone only waits for the container to start — not for the service inside to be ready. A `healthcheck` + `condition: service_healthy` makes the API wait until Postgres actually accepts connections.

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres"]
  interval: 5s      # check every 5 seconds
  timeout: 5s       # fail if no response in 5 seconds
  retries: 5        # mark unhealthy after 5 consecutive failures
```

---

## Phase 6 — Environment Variables, Secrets & Resource Limits

### Environment Variables

Config should never be hardcoded. Use an `env_file` in Compose:

```yaml
services:
  api:
    env_file:
      - .env
```

`.env` is loaded at `docker compose up` time. Never commit `.env` — commit `.env.example` with placeholder values.

### Docker Secrets

For sensitive values (passwords, tokens), secrets are more secure than env vars because they're mounted as files inside the container rather than exposed in `docker inspect` output.

```yaml
secrets:
  db_password:
    file: ./secrets/db_password.txt

services:
  postgres:
    secrets:
      - db_password
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
```

### Resource Limits

Prevent a single container from consuming all host resources:

```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: "0.5"      # max 50% of one CPU core
          memory: 256M     # max 256MB RAM
```

---

## Phase 7 — Image Tagging & Registry Push

### Tagging

```bash
docker build -t prachi-satbhai0741/taskflow-go:v1.0 .
docker build -t prachi-satbhai0741/taskflow-go:latest .
```

Use semantic versioning (`v1.0`, `v1.1`) for releases and `latest` for the most recent build.

### Pushing to DockerHub

```bash
docker login
docker push prachi-satbhai0741/taskflow-go:v1.0
docker push prachi-satbhai0741/taskflow-go:latest
```

Anyone can then pull and run your image:

```bash
docker pull prachi-satbhai0741/taskflow-go:latest
docker run -p 8080:8080 prachi-satbhai0741/taskflow-go:latest
```
