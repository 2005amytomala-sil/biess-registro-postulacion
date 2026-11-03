# Backend Setup

## Prerequisites
- Go 1.21+
- Docker and Docker Compose
- PostgreSQL client (optional)

## Running the application

### Start dependencies (PostgreSQL and Redis)
```bash
docker-compose up -d
```

### Load environment variables
```bash
cp .env.example .env
# Edit .env if needed (defaults match docker-compose)
```

### Run the server
```bash
go run ./cmd/api
```

Or use Makefile:
```bash
make run
```

## Health check
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "services": {
    "postgres": "healthy",
    "redis": "healthy"
  }
}
```

## Build
```bash
make build
```

## Stop dependencies
```bash
docker-compose down
```
