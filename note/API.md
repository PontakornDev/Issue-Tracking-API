# Issue Tracking API - Simplified

A RESTful API for managing issues built with **Go**, **Gin** framework, and **GORM** ORM with **PostgreSQL**.

## Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL 12+

### Setup
```bash
cd /Users/c.ptk/Desktop/product/Issue-Tracking

# Create database
psql -U postgres -c "CREATE DATABASE issue_tracking;"

# Run server
go run main.go
```

Server runs on `http://localhost:8080`

---

## API Endpoints

### 1. Create Issue
```
POST /api/issues
```

**Request Body:**
```json
{
  "reporter_id": 1,
  "assignee_id": 1,
  "status_id": 1,
  "title": "Login bug",
  "description": "Users cannot log in with special characters",
  "priority": "high"
}
```

**Response:** `201 Created` with issue object

---

### 2. Get Issues (with optional status filter)
```
GET /api/issues
GET /api/issues?status=open
```

**Query Parameters:**
- `status` (optional): Filter by status code (e.g., `open`, `in_progress`, `resolved`)

**Response:** `200 OK` with array of issues

---

### 3. Get Single Issue
```
GET /api/issues/:id
```

**Response:** `200 OK` with issue object including all relations

---

### 4. Update Issue Status
```
PATCH /api/issues/:id/status
```

**Request Body:**
```json
{
  "new_status_id": 2,
  "comment": "Assigned to John for review"
}
```

**Response:** `200 OK` with updated issue (status history is automatically recorded)

---

### 5. Create Comment on Issue
```
POST /api/issues/:id/comment
```

**Request Body:**
```json
{
  "user_id": 1,
  "content": "This is a comment on the issue"
}
```

**Response:** `201 Created` with comment object

---

## Example Requests

### Create an issue
```bash
curl -X POST http://localhost:8080/api/issues \
  -H "Content-Type: application/json" \
  -d '{
    "reporter_id": 1,
    "status_id": 1,
    "title": "Login bug",
    "description": "Special characters fail",
    "priority": "high"
  }'
```

### Get all open issues
```bash
curl http://localhost:8080/api/issues?status=open
```

### Get a specific issue with all details
```bash
curl http://localhost:8080/api/issues/1
```

### Update issue status and record history
```bash
curl -X PATCH http://localhost:8080/api/issues/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "new_status_id": 2,
    "comment": "Moved to in progress"
  }'
```

### Add a comment to an issue
```bash
curl -X POST http://localhost:8080/api/issues/1/comment \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "content": "Working on this now"
  }'
```

---

## Database Schema

**Tables:**
- `users` - Issue reporters
- `officer` - Issue handlers/assignees
- `issue_statuses` - Status definitions
- `issues` - Issue records
- `issue_status_history` - Status change tracking
- `comments` - Issue comments

---

## Response Format

All responses are JSON. Success responses follow this format:

```json
{
  "issue_id": 1,
  "title": "Login bug",
  "status": {...},
  "comments": [...]
}
```

Error responses:
```json
{
  "error": "Issue not found"
}
```

---

## Status Codes
- `200` - OK
- `201` - Created
- `400` - Bad Request
- `404` - Not Found
- `500` - Server Error
