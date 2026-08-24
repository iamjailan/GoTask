# GoTask

Simple CRUD API for tasks using Go, Gin, GORM, and PostgreSQL.

## Run locally

Make sure PostgreSQL is running locally and that the `gotask` database exists.
You can configure the connection in `.env` using the values from `.env.example`.

```bash
make setup
make run
```

The API listens on `http://localhost:8080`.

Run `make help` to see all available commands.

## Endpoints

- `GET /health`
- `POST /api/v1/tasks` — `{ "title": "Learn Go" }`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/:id`
- `PUT /api/v1/tasks/:id` — `{ "title": "Learn Gin", "completed": true }`
- `DELETE /api/v1/tasks/:id`
