# URL Shortener — Backend

A blazing-fast URL shortener REST API built with **Go**, **Gin**, **PostgreSQL**, and **Redis**.

## Features

- 🔗 **Shorten URLs** — auto-generated or custom short codes
- 👤 **User authentication** — JWT-based register/login
- 📊 **Analytics** — click tracking with device, OS, browser & referrer breakdowns
- ⚡ **Redis caching** — fast redirects with in-memory cache
- ⏳ **Link expiration** — set optional TTL on shortened URLs
- 🐳 **Docker support** — ready-to-build Dockerfile
- 🔒 **Protected routes** — per-user URL management

## Tech Stack

| Layer      | Technology                                  |
| ---------- | ------------------------------------------- |
| Language   | Go 1.21+                                    |
| Router     | [Gin](https://github.com/gin-gonic/gin)     |
| ORM        | [GORM](https://gorm.io/)                    |
| Database   | PostgreSQL                                  |
| Cache      | Redis (go-redis v9)                         |
| Auth       | JWT (golang-jwt v5)                         |
| Migration  | GORM AutoMigrate + manual SQL fallback      |

## Project Structure

```
backend/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Environment config loader
│   ├── database/
│   │   ├── db.go            # PostgreSQL connection
│   │   └── redis.go         # Redis connection
│   ├── handlers/
│   │   ├── analytics.go     # Stats endpoints
│   │   ├── auth.go          # Register / login / logout
│   │   └── url.go           # URL CRUD + redirect
│   ├── middleware/
│   │   └── auth.go          # JWT auth middleware
│   ├── models/
│   │   ├── analytics.go     # Analytics models
│   │   ├── url.go           # URL models
│   │   └── user.go          # User models
│   ├── repository/
│   │   ├── url_repo.go      # URL DB + cache operations
│   │   └── user_repo.go     # User DB operations
│   └── utils/
│       ├── jwt.go           # Token generation & validation
│       ├── password.go      # bcrypt hashing
│       └── random.go        # Short code generator
├── migrations/
│   └── init.sql             # Manual SQL migration (optional)
├── Dockerfile
├── .env.example
├── go.mod
├── go.sum
└── .env.example
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis
- (Optional) Docker & Docker Compose

### Environment Variables

Copy the example and adjust values:

```bash
cp .env.example .env
```

| Variable          | Default                                 | Description                |
| ----------------- | --------------------------------------- | -------------------------- |
| `PORT`            | `8080`                                  | Server port                |
| `DB_HOST`         | `localhost`                             | PostgreSQL host            |
| `DB_USER`         | `postgres`                              | PostgreSQL user            |
| `DB_PASSWORD`     | `postgres`                              | PostgreSQL password        |
| `DB_NAME`         | `urlshortener`                          | Database name              |
| `DB_PORT`         | `5432`                                  | PostgreSQL port            |
| `DB_SSL_MODE`     | `disable`                               | SSL mode                   |
| `REDIS_HOST`      | `localhost`                             | Redis host                 |
| `REDIS_PORT`      | `6379`                                  | Redis port                 |
| `REDIS_PASSWORD`  | _(empty)_                               | Redis password             |
| `JWT_SECRET`      | `your-super-secret-jwt-key-change-this` | JWT signing secret         |
| `JWT_EXPIRY_HOURS`| `24`                                    | Token expiry in hours      |
| `BASE_URL`        | `http://localhost:8080`                 | Public base URL for links  |

### Run Locally

```bash
# Install dependencies
go mod download

# Start PostgreSQL & Redis (if not running)
# ...

# Run the server
go run ./cmd/main.go
```

Server starts on `http://localhost:8080`.

### Docker

```bash
# Build image
docker build -t url-shortener .

# Run container
docker run -p 8080:8080 --env-file .env url-shortener
```

## API Reference

### Public

| Method | Endpoint              | Description              |
| ------ | --------------------- | ------------------------ |
| `GET`  | `/health`             | Health check             |
| `GET`  | `/:code`              | Redirect to original URL |
| `GET`  | `/api/urls/:code/stats` | Public click stats     |

### Auth

| Method | Endpoint             | Body                                     | Description     |
| ------ | -------------------- | ---------------------------------------- | --------------- |
| `POST` | `/api/auth/register` | `{ email, password, name }`             | Create account  |
| `POST` | `/api/auth/login`    | `{ email, password }`                   | Get JWT token   |
| `POST` | `/api/auth/logout`   | —                                        | Logout          |

### Protected (Bearer token required)

| Method   | Endpoint             | Body                                       | Description      |
| -------- | -------------------- | ------------------------------------------ | ---------------- |
| `POST`   | `/api/urls`          | `{ original_url, custom_code?, expires_at? }` | Shorten URL   |
| `GET`    | `/api/urls`          | —                                          | List user's URLs |
| `GET`    | `/api/urls/:code`    | —                                          | Get URL details  |
| `PUT`    | `/api/urls/:code`    | `{ original_url?, is_active?, expires_at? }` | Update URL    |
| `DELETE` | `/api/urls/:code`    | —                                          | Delete URL       |

### Example Requests

**Register:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123","name":"John"}'
```

**Create short URL:**
```bash
curl -X POST http://localhost:8080/api/urls \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://example.com/very-long-url"}'
```

**Response:**
```json
{
  "id": 1,
  "original_url": "https://example.com/very-long-url",
  "short_url": "http://localhost:8080/aB3xY7",
  "short_code": "aB3xY7",
  "clicks": 0,
  "is_active": true,
  "created_at": "2025-01-01T00:00:00Z"
}
```

## Database

GORM auto-migrates the schema on startup. For manual setup, run:

```bash
psql -U postgres -d urlshortener -f migrations/init.sql
```

### Tables

- **users** — registered users
- **urls** — shortened URLs with expiry & click counts
- **analytics** — per-click data (device, OS, browser, referrer, geo)

## License

MIT
