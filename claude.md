# CLAUDE.md — valeth-urls-shortener

## Project Overview

A URL shortener REST API built in Go using Domain-Driven Design (DDD).
Only registered users can create short URLs. Built as a learning and portfolio project.

**Stack:** Go, PostgreSQL, Fiber v2, GORM, JWT authentication

---

## Architecture — DDD Layer Structure

cmd/api/main.go → entry point, calls bootstrap only
bootstrap/bootstrap.go → wires all dependencies (composition root)
config/config.go → loads env vars, exports Config struct
internal/
domain/url/ → URL entity, repository interface, domain errors
domain/user/ → User entity, repository interface
app/url/ → URL use cases (CreateURL, DeleteURL, ListUserURLs)
app/auth/ → Auth use cases (Register, Login)
infra/persistence/ → GORM implementations of repository interfaces
infra/database/ → DB connection and AutoMigrate
interfaces/http/ → Fiber handlers, router, JWT middleware
dto/ → Request and response structs with validate tags

### Layer Rules — Strictly Enforced

- **Domain layer** — pure Go only. No Fiber, no GORM imports. Contains entities and repository interfaces.
- **Application layer** — orchestrates domain. No Fiber imports. Returns DTOs, never raw entities.
- **Infrastructure layer** — GORM lives here only. Implements repository interfaces.
- **Interface layer** — Fiber handlers only. Parses HTTP → calls application layer → writes response. Zero business logic.
- **Bootstrap** — the ONLY place that wires dependencies together. Nothing outside bootstrap reaches app.DB directly.

---

## Key Design Decisions

### Entity + GORM Model Combined

Domain entities use GORM tags directly (pragmatic tradeoff for this project scale).
Do NOT suggest separating them into domain structs + GORM models + mappers.

### Bootstrap Pattern

`App` struct holds `*fiber.App`, `*gorm.DB`, `*config.Config`.
`main.go` only calls `bootstrap.NewApp()` and `app.Fiber.Listen()`.

### Short Code

- User always provides a custom short code (no auto-generation in current scope)
- Validated as `alphanum`, `min=4`, `max=30`
- Uniqueness checked before insert

### Short URL Construction

Never stored in DB. Always constructed at response time:

```go
shortURL = cfg.BaseURL + "/" + shortCode
```

### Delete Endpoint

Uses path param only — no request body.

DELETE /api/v1/urls/:code ✓
DELETE /api/v1/urls (body) ✗

### Soft Delete

URLs use `DeletedAt *time.Time` for soft delete, not hard delete.

---

## Naming Conventions

- Go acronyms: `OriginalURL`, `ShortCode`, `UserID` — never `OriginalUrl`, `UserId`
- JSON tags: snake_case — `original_url`, `short_code`, `created_at`
- Repository methods: `Create`, `FindByShortCode`, `FindByUserID`, `Delete`

---

## Entities

### URL

```go
type URL struct {
    ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    UserID      uuid.UUID  `gorm:"type:uuid;not null"`
    OriginalURL string     `gorm:"type:text;not null"`
    ShortCode   string     `gorm:"type:varchar(30);uniqueIndex;not null"`
    CreatedAt   time.Time  `gorm:"not null"`
    UpdatedAt   time.Time  `gorm:"not null"`
    DeletedAt   *time.Time `gorm:"index"`
}
```

### User

```go
type User struct {
    ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Email        string     `gorm:"type:varchar(255);uniqueIndex;not null"`
    PasswordHash string     `gorm:"type:text;not null"`
    CreatedAt    time.Time  `gorm:"not null"`
    UpdatedAt    time.Time  `gorm:"not null"`
}
```

---

## API Endpoints

| Method | Path                 | Auth | Description              |
| ------ | -------------------- | ---- | ------------------------ |
| POST   | `/auth/register`     | No   | Register new user        |
| POST   | `/auth/login`        | No   | Login, returns JWT       |
| POST   | `/api/v1/urls`       | JWT  | Create short URL         |
| GET    | `/api/v1/urls`       | JWT  | List user's URLs         |
| DELETE | `/api/v1/urls/:code` | JWT  | Delete a URL             |
| GET    | `/:code`             | No   | Redirect to original URL |

---

## Response Shape

```json
{
  "id": "uuid",
  "original_url": "https://example.com/long/path",
  "short_code": "mylink",
  "short_url": "https://short.ly/mylink",
  "created_at": "2026-04-27T10:00:00Z",
  "updated_at": "2026-04-27T10:00:00Z"
}
```

---

## Out of Scope — Do Not Suggest

- Click analytics / click count tracking
- Link expiry / `expires_at`
- Anonymous URL creation
- Auto-generated short codes (user always provides the code)
- Separate GORM models from domain entities

---

## Local Development

```bash
# Run the app
go run cmd/api/main.go

# Database is PostgreSQL running in OrbStack Docker container
# Container name: urlshortener-postgres
# Connection via DATABASE_URL in .env

# .env required fields
DATABASE_URL=postgres://admin:secret@localhost:5432/urlshortener
BASE_URL=http://localhost:8080
JWT_SECRET=your-secret-here
```

---

## Module Path

`github.com/haramurti/valeth-urls-shortener`
