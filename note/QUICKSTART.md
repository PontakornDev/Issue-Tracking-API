# Quick Start Guide

## Run with Docker Compose (Recommended)

### Prerequisites
- Docker & Docker Compose installed

### Start Services
```bash
cd /Users/c.ptk/Desktop/product/Issue-Tracking

# Build and start services
docker-compose up --build

# Or start in background
docker-compose up -d --build
```

### Check Logs
```bash
# All services
docker-compose logs -f

# Only API
docker-compose logs -f api

# Only Database
docker-compose logs -f postgres
```

### Stop Services
```bash
docker-compose down

# Also remove volumes
docker-compose down -v
```

---

## Run Locally (Without Docker)

### Prerequisites
- Go 1.24+
- PostgreSQL running locally

### Setup Database
```bash
psql -U postgres -c "CREATE DATABASE issue_tracking;"
```

### Run Server
```bash
cd /Users/c.ptk/Desktop/product/Issue-Tracking
go run main.go
```

Server runs on `http://localhost:8080`

---

## Environment Variables

### Docker (auto-configured)
- `DATABASE_URL`: `host=postgres user=postgres password=postgres dbname=issue_tracking port=5432 sslmode=disable`

### Local
```bash
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=issue_tracking port=5432 sslmode=disable"
go run main.go
```

---

## Service Ports

| Service | Port | URL |
|---------|------|-----|
| PostgreSQL | 5432 | `localhost:5432` |
| API | 8080 | `http://localhost:8080` |

---

## API Health Check

```bash
curl http://localhost:8080/health
```

Response: `{"status":"ok"}`

---

## Docker Service Architecture

```
┌─────────────────────────────────────┐
│   Docker Network                    │
│  (issue-tracking-network)           │
├─────────────┬───────────────────────┤
│             │                       │
│  postgres   │      api              │
│  :5432      │      :8080            │
│             │                       │
│ (database)  │   (Go + Gin + GORM)   │
└─────────────┴───────────────────────┘
```

**Database Service (postgres):**
- Alpine Linux PostgreSQL 16
- Persistent volume: `postgres_data`
- Health check enabled
- Auto-creates database on startup

**API Service (api):**
- Multi-stage Docker build (optimized)
- Depends on postgres service
- Environment variables configured
- Port 8080 exposed

---

## Troubleshooting

### Port Already in Use
```bash
# Kill process on port 5432
lsof -ti:5432 | xargs kill -9

# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

### Database Connection Error
```bash
# Check if postgres is running
docker-compose logs postgres

# Restart postgres
docker-compose restart postgres
```

### Rebuild Everything
```bash
docker-compose down -v
docker-compose up --build
```
