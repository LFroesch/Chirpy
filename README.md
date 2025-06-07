# Chirpy

A backend only twitter-like social media API built with Go, PostgreSQL, and JWT authentication.
For v2 (with frontend) checkout: https://github.com/LFroesch/social-media-demo

## Features

- User registration and authentication
- Create, read, and delete chirps (140 character limit)
- JWT access tokens with refresh token support
- Profanity filtering
- Premium user upgrades via webhook
- PostgreSQL database with migrations

## Quick Start

### Prerequisites
- Go 1.22.6+
- PostgreSQL
- Goose (for migrations)

### Setup
1. Clone the repository
2. Copy `.env.example` to `.env` and configure:
   ```
   DB_URL=postgres://postgres:postgres@localhost:5432/chirpy
   PLATFORM=dev
   JWT_SECRET=your-secret-key
   POLKA_KEY=your-webhook-key
   ```
3. Run database migrations:
   ```bash
   make migrate-up
   ```
4. Start the server:
   ```bash
   go run .
   ```

The server runs on `http://localhost:8080`

## API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/healthz` | Health check | No |
| POST | `/api/users` | Create user | No |
| POST | `/api/login` | User login | No |
| POST | `/api/refresh` | Refresh JWT token | Yes |
| POST | `/api/revoke` | Revoke refresh token | Yes |
| PUT | `/api/users` | Update user | Yes |
| POST | `/api/chirps` | Create chirp | Yes |
| GET | `/api/chirps` | List chirps | No |
| GET | `/api/chirps/{id}` | Get chirp by ID | No |
| DELETE | `/api/chirps/{id}` | Delete chirp | Yes |
| POST | `/api/polka/webhooks` | Upgrade webhook | API Key |
| GET | `/admin/metrics` | View metrics | No |
| POST | `/admin/reset` | Reset database (dev only) | No |

## Authentication

- Access tokens expire in 1 hour
- Refresh tokens expire in 60 days
- Include JWT in Authorization header: `Bearer <token>`
- API webhook requires: `Authorization: ApiKey <key>`

## Database Commands

```bash
# Run migrations
make migrate-up

# Rollback migrations  
make migrate-down

# Reset database
make reset-db

# Connect to PostgreSQL
make psql
```

## Development

The project uses:
- **SQLC** for type-safe SQL queries
- **Goose** for database migrations  
- **JWT** for authentication
- **Bcrypt** for password hashing