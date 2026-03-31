# my-note — Go + Gin in-memory notes API

Lightweight development API for notes (save and retrieve). No auth, in-memory storage.

Run locally:

```bash
go run main.go
```

Build and run with Docker Compose:

```bash
docker compose up --build
```

API endpoints:

- `POST /api/notes` — create a note. JSON: `{ "title": "...", "content": "..." }`
- `GET /api/notes` — list all notes
- `GET /api/notes/:id` — get a single note

Example curl:

```bash
curl -X POST localhost:8080/api/notes -H "Content-Type: application/json" \
  -d '{"title":"Hello","content":"This is a note"}'

curl localhost:8080/api/notes
curl localhost:8080/api/notes/<id>
```

Notes:

- Data is stored in-memory and will be lost when the process stops.
- Run `go mod tidy` to populate dependencies before `go build` or `go run` if needed.
