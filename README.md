# Workout Scheduler (PoC)

This is a small proof-of-concept workout scheduling app in Go.

Features
- Create accounts (`/signup`)
- Log in and receive a JWT (`/login`)
- List and create classes (`/classes`)
- Sign up for a class (`/classes/{id}/signup`)
- Log activities and list them (`/activities`)

Quick start

1. Install dependencies and sqlc (optional to generate types):

   ```sh
   go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
   go mod tidy
   ```

2. (Optional) Generate db types with `sqlc`:

   ```sh
   sqlc generate
   ```

3. Run the server:

   ```sh
   go run ./cmd/server
   ```

API examples

- Sign up:
  POST /signup {"email":"a@a.com","password":"p","full_name":"A"}
- Login:
  POST /login {"email":"a@a.com","password":"p"} -> returns {"token":"..."}
- List classes:
  GET /classes
- Sign up for class (auth header `Authorization: Bearer <token>`):
  POST /classes/1/signup
- Log activity:
  POST /activities {"type":"run","duration_minutes":30,"notes":"nice"}

Notes
-- The project includes `sqlc.yaml` and SQL files under `internal/db/migrations` and `internal/db/queries` for use with `sqlc`.
- The app uses an in-memory SQLite DB; restarting the server clears data.
# test-app-2
test application for ci/cd tools, v2
