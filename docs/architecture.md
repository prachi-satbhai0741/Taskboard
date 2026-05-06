# Architecture

## Overview

`taskflow-go` is a three-tier containerized application — a Go REST API, a PostgreSQL database for persistent task storage, and a Redis cache for future rate limiting and session support. All three services run inside an isolated Docker network.

The architecture is deliberately kept simple. The goal is not a complex distributed system — it is a clean, real-world containerization pattern that demonstrates how services are connected, isolated, and orchestrated using Docker.

---

## Container Diagram

```
┌─────────────────────────────────────────────────────┐
│                Docker Network: taskflow-net         │
│                                                     │
│   ┌─────────────────┐       ┌────────────────────┐  │
│   │   taskflow-api  │──────▶│  taskflow-postgres │  │
│   │   Go + Gin      │       │  PostgreSQL 15     │  │
│   │   :8080         │       │  :5432             │  │
│   └────────┬────────┘       └────────────────────┘  │
│            │                 Volume: pg-data        │
│            │                ┌────────────────────┐  │
│            └───────────────▶│  taskflow-redis    │  │
│                             │  Redis 7           │  │
│                             │  :6379             │  │
│                             └────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

**Port exposed to host:** only `:8080` (API). Postgres and Redis are internal — not reachable from outside the Docker network.

---

## Tech Stack

| Layer | Technology | Reason |
|---|---|---|
| Language | Go 1.22 | Compiles to a single static binary — ideal for minimal Docker images |
| Framework | Gin | Lightweight, fast HTTP router; minimal boilerplate |
| ORM | GORM | Handles migrations and queries cleanly without raw SQL |
| Database | PostgreSQL 15 Alpine | Relational, production-standard, lightweight Alpine image |
| Cache | Redis 7 Alpine | In-memory store; demonstrates multi-service Compose setup |
| Container | Docker + Compose | Core focus of the project |

---

## Why Go for Docker Learning?

Go compiles to a single static binary with no runtime dependencies. This makes it the best language to demonstrate multi-stage Docker builds — the final image can be as small as ~20MB using `scratch` or `distroless`, compared to 300MB+ for Node.js or Python equivalents.

```
Stage 1 (golang:1.22-alpine)  →  builds binary
Stage 2 (scratch or alpine)   →  copies binary only
                               →  no compiler, no source, no dependencies
```

---

## Data Flow

```
Client Request
     │
     ▼
  :8080 (host)
     │
     ▼
taskflow-api container (Gin router)
     │
     ├──▶ POST/GET/PUT/DELETE /tasks ──▶ GORM ──▶ PostgreSQL
     │
     └──▶ GET /health ──▶ 200 OK
```

---

## Design Decisions

**Single binary, no CGO** — `CGO_ENABLED=0` is set during build so the binary is fully static and runs in a `scratch` base image.

**Alpine for Postgres and Redis** — Alpine-based official images are used to keep the overall stack size minimal and consistent with production best practices.

**Internal-only DB ports** — Postgres and Redis ports are not published to the host in production Compose. They are only accessible within `taskflow-net`. Dev overrides in `docker-compose.override.yml` expose them for local inspection.

**Health check before app start** — the API container depends on Postgres with a health check condition, so GORM never tries to connect to a database that isn't ready yet.
