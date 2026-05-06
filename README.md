<div align="center">

# Taskboard

**A production-style Developer Task Board REST API built in Go**
*A hands-on Docker consolidation project covering every core Docker concept — from basic Dockerfiles to multi-stage builds, networking, Compose, secrets, and registry publishing.*

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=flat-square&logo=docker&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat-square&logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

</div>

---

## What is this?

`Taskboard` is a minimal REST API for managing developer tasks — think a Jira backend, kept intentionally simple so the **focus stays entirely on Docker**. Every phase of the project introduces a new Docker concept cluster in a real, working context.

> Each phase = one commit = one Docker concept group. The git history tells the full story.

---

## Quick Start

```bash
git clone https://github.com/prachi-satbhai0741/Taskboard.git
cd Taskboard
cp .env.example .env
docker compose up --build
```

API available at `http://localhost:8080` — see [`docs/api-reference.md`](docs/api-reference.md) for all endpoints.

---

## Project Structure

```
Taskboard/
├── Dockerfile                   # Multi-stage production build
├── docker-compose.yml           # Orchestrates API + Postgres + Redis
├── docker-compose.override.yml  # Dev overrides
├── .dockerignore
├── .env.example                 # Copy to .env and fill values
├── main.go
├── go.mod / go.sum
├── handlers/                    # HTTP route handlers
├── models/                      # GORM models
├── db/                          # DB connection + migration
├── secrets/                     # Docker secrets (git-ignored)
└── docs/
    ├── architecture.md          # System design + diagrams
    ├── docker-concepts.md       # All Docker concepts used + explanations
    ├── setup-guide.md           # Full local + Docker setup instructions
    ├── api-reference.md         # All endpoints with curl examples
    └── project-phases.md        # Phase-by-phase breakdown with git log
```

---

## Docker Concepts Covered

| Concept | Phase |
|---|---|
| Dockerfile, `.dockerignore` | 2 |
| Multi-stage build, layer caching | 3 |
| Custom bridge networks | 4 |
| Volumes, Docker Compose, health checks | 5 |
| Env vars, secrets, resource limits | 6 |
| Image tagging, DockerHub push | 7 |

Full breakdown → [`docs/docker-concepts.md`](docs/docker-concepts.md)

---

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Service health check |
| `GET` | `/tasks` | List all tasks |
| `POST` | `/tasks` | Create a new task |
| `PUT` | `/tasks/:id` | Update task status |
| `DELETE` | `/tasks/:id` | Delete a task |

Full reference with examples → [`docs/api-reference.md`](docs/api-reference.md)

---

## Phases

| Phase | Focus | Docker Concepts |
|---|---|---|
| **1** | Go REST API baseline — no Docker | None — app foundation |
| **2** | Containerize the app | Dockerfile, `.dockerignore` |
| **3** | Optimize the image | Multi-stage build, layer caching |
| **4** | Connect services | Custom bridge networks |
| **5** | Orchestrate everything | Docker Compose, volumes, health checks |
| **6** | Harden config | Env vars, secrets, resource limits |
| **7** | Ship it | Image tagging, DockerHub push |

Full breakdown with commands → [`docs/project-phases.md`](docs/project-phases.md)

---

## Docs

| File | Contents |
|---|---|
| [`architecture.md`](docs/architecture.md) | System design, container diagram, tech choices |
| [`docker-concepts.md`](docs/docker-concepts.md) | Every Docker concept used with explanation |
| [`setup-guide.md`](docs/setup-guide.md) | Local + Docker + Compose setup instructions |
| [`api-reference.md`](docs/api-reference.md) | Endpoints, request/response, curl examples |
| [`project-phases.md`](docs/project-phases.md) | Phase breakdown, commit messages, what was learned |

---

## Related Repos

| Repo | Stack | What it covers |
|---|---|---|
| [`student-api-go`](https://github.com/prachi-satbhai0741/student-api-go) | Go + Gin | Go REST API fundamentals |
| [`Docker-bridge-url-shortener`](https://github.com/prachi-satbhai0741/Docker-bridge-url-shortener) | Node.js + Redis | REST API + Redis caching |
| [`Docker_Lab`](https://github.com/prachi-satbhai0741/Docker_Lab) | Docker | Module-wise Docker hands-on exercises |

---

## Author

**Prachi Satbhai**
[![GitHub](https://img.shields.io/badge/GitHub-prachi--satbhai0741-181717?style=flat-square&logo=github)](https://github.com/prachi-satbhai0741)

---

## License

This project is licensed under the MIT License — see LICENSE for details.[`LICENSE`](LICENSE)

---

If this repo helped you, consider giving it a star — it helps others find it too!
