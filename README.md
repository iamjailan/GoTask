# GoTask

Simple CRUD API for tasks using Go, Gin, GORM, and PostgreSQL.

## Run locally

Make sure PostgreSQL is running locally and that the `gotask` database exists.
You can configure the connection in `.env` using the values from `.env.example`.

```bash
make setup
make run
```

Schema changes are tracked as ordered SQL files in `migrations/`. Run `make migration`, enter a descriptive name such as `add due date to tasks`, and it creates matching files such as `20260827143000_add_due_date_to_tasks.up.sql` and `.down.sql`. Apply pending migrations with `make migrate`; `make migrate-down` rolls back one migration, and `make migrate-version` shows the recorded version. Use `make migrate-reset` to drop and recreate PostgreSQL's `public` schema. This deletes all tables, data, indexes, sequences, and migration records in that schema; run `make migrate` afterward to rebuild it from the migration files. Each `make migrate` run creates or updates the `migration_history` table, which records `version`, `name`, `created_at`, `applied_at`, and `is_applied` for every migration file. The API does not change the schema during startup.

`make run` uses [Air](https://github.com/air-verse/air) for hot reload: changes to Go files automatically rebuild and restart the API. Install it once if needed:

```bash
go install github.com/air-verse/air@v1.61.0
```

To run the API once without a file watcher, use `make start`.

The API listens on `http://localhost:8080`.

Run `make help` to see all available commands.

## Endpoints

- `GET /health`
- `GET /swagger/index.html` — Swagger UI (requires HTTP Basic Auth)
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

## Swagger documentation

Open `/swagger/index.html` at the configured API address. The UI, OpenAPI JSON, and Swagger assets require HTTP Basic Auth. Set `SWAGGER_USERNAME` and `SWAGGER_PASSWORD` in `.env`; the API will not start without both values.

Run `make swagger` after changing API annotations to regenerate the checked-in `docs/swagger.json` and `docs/swagger.yaml` files.

## Rate limits

Customer API requests are limited to 10 requests per minute per client IP. Email-sending endpoints (`POST /customer/auth/register` and `PUT /customer/me/email`) are limited to 3 requests per minute per client IP. Exceeding a limit returns `429 Too Many Requests` with a `Retry-After` header.

Set `JWT_SECRET` in `.env` to a long random value before running outside local development.
