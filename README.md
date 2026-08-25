# GoTask

Simple CRUD API for tasks using Go, Gin, GORM, and PostgreSQL.

## Run locally

Make sure PostgreSQL is running locally and that the `gotask` database exists.
You can configure the connection in `.env` using the values from `.env.example`.

```bash
make setup
make run
```

`make run` uses [Air](https://github.com/air-verse/air) for hot reload: changes to Go files automatically rebuild and restart the API. Install it once if needed:

```bash
go install github.com/air-verse/air@v1.61.0
```

To run the API once without a file watcher, use `make start`.

The API listens on `http://localhost:8080`.

Run `make help` to see all available commands.

## Endpoints

- `GET /health`
- `POST /api/v1/tasks` — `{ "title": "Learn Go", "description": "...", "status": "pending", "priority": "medium", "due_date": "2026-09-01T12:00:00Z" }` (requires `Authorization: Bearer <token>`)
- `GET /api/v1/tasks` (requires `Authorization: Bearer <token>`)
- `GET /api/v1/tasks/:id` — IDs use the `tsk_<UUID>` format
- `PUT /api/v1/tasks/:id` — `{ "title": "Learn Gin", "completed": true }`
- `DELETE /api/v1/tasks/:id`
- `POST /api/v1/auth/register` — create a customer and return a 24-hour JWT
- `POST /api/v1/auth/login` — authenticate a customer and return a 24-hour JWT

Set `JWT_SECRET` in `.env` to a long random value before running outside local development.
