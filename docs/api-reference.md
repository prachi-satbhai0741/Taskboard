# API Reference

Base URL: `http://localhost:8080`

---

## Task Status Values

```
todo | in_progress | done
```

---

## Endpoints

### `GET /health`

Returns service health status.

**Response**
```json
{
  "status": "ok",
  "service": "taskflow-go"
}
```

---

### `GET /tasks`

Returns all tasks (soft-deleted tasks excluded).

**Response**
```json
[
  {
    "id": 1,
    "title": "Write Dockerfile",
    "description": "Multi-stage build for Go app",
    "status": "in_progress",
    "created_at": "2026-05-06T10:00:00Z",
    "updated_at": "2026-05-06T11:00:00Z"
  }
]
```

**curl**
```bash
curl http://localhost:8080/tasks
```

---

### `POST /tasks`

Creates a new task. Status defaults to `todo` if not provided.

**Request Body**
```json
{
  "title": "Write Dockerfile",
  "description": "Multi-stage build for Go app",
  "status": "todo"
}
```

**Response** `201 Created`
```json
{
  "id": 1,
  "title": "Write Dockerfile",
  "description": "Multi-stage build for Go app",
  "status": "todo",
  "created_at": "2026-05-06T10:00:00Z",
  "updated_at": "2026-05-06T10:00:00Z"
}
```

**curl**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Write Dockerfile", "description": "Multi-stage build for Go app"}'
```

---

### `PUT /tasks/:id`

Updates the status of an existing task.

**URL Params**
- `id` — task ID (integer)

**Request Body**
```json
{
  "status": "done"
}
```

**Response** `200 OK`
```json
{
  "id": 1,
  "title": "Write Dockerfile",
  "description": "Multi-stage build for Go app",
  "status": "done",
  "created_at": "2026-05-06T10:00:00Z",
  "updated_at": "2026-05-06T12:00:00Z"
}
```

**curl**
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "done"}'
```

**Error — Task Not Found** `404`
```json
{
  "error": "Task not found"
}
```

---

### `DELETE /tasks/:id`

Soft-deletes a task (sets `deleted_at`, not permanently removed from DB).

**URL Params**
- `id` — task ID (integer)

**Response** `200 OK`
```json
{
  "message": "Task deleted"
}
```

**curl**
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

**Error — Task Not Found** `404`
```json
{
  "error": "Task not found"
}
```

---

## Error Format

All errors follow this structure:

```json
{
  "error": "description of what went wrong"
}
```

| Status Code | Meaning |
|---|---|
| `200` | OK |
| `201` | Created |
| `400` | Bad request — invalid JSON or missing required field |
| `404` | Task not found |
| `500` | Internal server error |
