# GoTask

Simple CRUD API for tasks using Go, Gin, GORM, and PostgreSQL.

## Run locally

Make sure PostgreSQL is running locally and that the `gotask` database exists.
You can configure the connection in `.env` using the values from `.env.example`.

```bash
make setup
make run
```

Schema changes are tracked as ordered SQL files in `migrations/`. Apply pending migrations with `make migrate`; `make migrate-down` rolls back one migration, and `make migrate-version` shows the recorded version. Create a new numbered `*.up.sql` and matching `*.down.sql` file for every schema change. The API does not change the schema during startup.

`make run` uses [Air](https://github.com/air-verse/air) for hot reload: changes to Go files automatically rebuild and restart the API. Install it once if needed:

```bash
go install github.com/air-verse/air@v1.61.0
```

To run the API once without a file watcher, use `make start`.

The API listens on `http://localhost:8080`.

Run `make help` to see all available commands.

## Endpoints

- `GET /health`
- `POST /customer/tasks` — `{ "title": "Learn Go", "description": "...", "status": "pending", "priority": "medium", "due_date": "2026-09-01T12:00:00Z" }` (requires `Authorization: Bearer <token>`)
- `GET /customer/tasks` (requires `Authorization: Bearer <token>`)
- `GET /customer/tasks/:id` — IDs use the `tsk_<UUID>` format
- `PUT /customer/tasks/:id` — `{ "title": "Learn Gin", "completed": true }`
- `DELETE /customer/tasks/:id`
- `POST /customer/auth/register` — create a customer and return a 24-hour JWT
- `POST /customer/auth/register` — create a pending registration and send a 6-digit email verification code (`202 Accepted`; no customer account exists yet)
- `POST /customer/auth/confirm-email` — `{ "email": "user@example.com", "code": "123456" }`; verify the code and return a 24-hour JWT
- `POST /customer/auth/login` — authenticate a customer and return a 24-hour JWT
- `GET /customer/me` — get the authenticated user
- `PUT /customer/me` — update the authenticated user's profile (`first_name`, `last_name`, `phone`, `avatar_url`)
- `PUT /customer/me/email` — change the email with `{ "email": "new@example.com", "current_password": "..." }`; sends a change notification to the previous email address
- `PUT /customer/me/password` — change the password with `{ "current_password": "...", "new_password": "..." }`
- `DELETE /customer/me` — delete the authenticated user

Set `JWT_SECRET` in `.env` to a long random value before running outside local development.
