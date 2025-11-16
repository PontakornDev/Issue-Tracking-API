# Issue Tracking API

A RESTful API for managing issues/support tickets built with **Go**, **Gin** framework, and **GORM** ORM with **PostgreSQL**.

## Features

- ✅ Full CRUD operations for issues
- ✅ PostgreSQL database with auto-migrations
- ✅ RESTful API endpoints
- ✅ JSON request/response handling
- ✅ Health check endpoint

## Prerequisites

- Go 1.23+ (project uses Go 1.24.10)
- PostgreSQL 12+ installed and running
- psql CLI (optional, for database management)

## Setup & Installation

### 1. Create PostgreSQL database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database and user (run in psql)
CREATE DATABASE issue_tracking;
CREATE USER postgres WITH PASSWORD 'postgres';
ALTER ROLE postgres SET client_encoding TO 'utf8';
ALTER ROLE postgres SET default_transaction_isolation TO 'read committed';
ALTER ROLE postgres SET default_transaction_deferrable TO on;
ALTER ROLE postgres SET default_transaction_read_only TO off;
GRANT ALL PRIVILEGES ON DATABASE issue_tracking TO postgres;
\q
```

Or use environment variable for custom connection:
```bash
export DATABASE_URL="host=localhost user=youruser password=yourpassword dbname=issue_tracking port=5432 sslmode=disable"
```

### 2. Install dependencies
```bash
cd /Users/c.ptk/Desktop/product/Issue-Tracking
go mod tidy
```

### 3. Run the server
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Building for Production

```bash
go build -o issue-tracking main.go
./issue-tracking
```

## API Endpoints

### Health Check
```
GET /health
```

### Issues
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/issues` | List all issues |
| GET | `/api/issues/:id` | Get a single issue |
| POST | `/api/issues` | Create a new issue |
| PUT | `/api/issues/:id` | Update an issue |
| DELETE | `/api/issues/:id` | Delete an issue |

## Example Requests

### Create an Issue
```bash
curl -X POST http://localhost:8080/api/issues \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Login bug",
    "description": "Users cannot log in with special characters",
    "status": "open"
  }'
```

### Get All Issues
```bash
curl http://localhost:8080/api/issues
```

### Update an Issue
```bash
curl -X PUT http://localhost:8080/api/issues/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Login bug",
    "description": "Users cannot log in with special characters",
    "status": "in-progress"
  }'
```

### Delete an Issue
```bash
curl -X DELETE http://localhost:8080/api/issues/1
```

## Project Structure

```
Issue-Tracking/
├── main.go          # Main application and API handlers
├── go.mod           # Module file with dependencies
├── go.sum           # Checksums for dependencies
├── .gitignore       # Git ignore rules
├── README.md        # This file
└── issues.db        # SQLite database (created on first run)
```

## Database

The application uses PostgreSQL with GORM for ORM. Connection details:
- **Default Host:** `localhost`
- **Default Port:** `5432`
- **Default User:** `postgres`
- **Default Password:** `postgres`
- **Default Database:** `issue_tracking`

Override via `DATABASE_URL` environment variable:
```bash
export DATABASE_URL="host=your-host user=your-user password=your-pass dbname=your-db port=5432 sslmode=disable"
go run main.go
```

### Issue Schema
```go
type Issue struct {
    ID          uint   // Primary key
    Title       string // Issue title
    Description string // Issue description
    Status      string // Status: open, in-progress, closed
    CreatedAt   int64  // Creation timestamp
    UpdatedAt   int64  // Last update timestamp
}
```

## Development Notes

- To add new models, define structs in `main.go` and call `db.AutoMigrate(&NewModel{})`
- Gin automatically handles JSON marshaling/unmarshaling
- PostgreSQL is production-ready and recommended for scalability
- Set `DATABASE_URL` environment variable for easy deployment configuration

## License

MIT
