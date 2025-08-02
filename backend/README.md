# Stock Insights Backend

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- PostgreSQL/CockroachDB
- Environment variables configured

### Installation
```bash
# Install dependencies
make deps

# Build the application
make build
```

## 📦 Components

### 1. API Server
The main API server that handles HTTP requests.

```bash
# Run API server
make run-api

# Or directly
go run cmd/api/main.go
```

**Endpoints:**
- `POST /api/admin/ingest` - Manual ingestion trigger
- `GET /api/admin/jobs/:jobId` - Job status tracking
- `GET /health` - Health check

### 2. Scheduled Worker
Automatically runs ingestion every 24 hours.

```bash
# Run scheduler worker
make run-scheduler

# Or directly
go run cmd/worker/scheduler/main.go
```

**Features:**
- ✅ Runs every 24 hours automatically
- ✅ Uses JobManager for tracking
- ✅ Graceful shutdown handling
- ✅ Error handling and retries

## 🔄 Architecture

### Manual Ingestion
```
POST /api/admin/ingest
    ↓
JobManager.CreateJob()
    ↓
JobManager.RunJobAsync()
    ↓
TriggerIngestionAsync()
    ↓
FetchAndProcessStocks()
```

### Scheduled Ingestion
```
Scheduler (24h)
    ↓
JobManager.CreateJob()
    ↓
JobManager.RunJobAsync()
    ↓
TriggerIngestionAsync()
    ↓
FetchAndProcessStocks()
```

## 🛠️ Development

### Run Both Services
```bash
# Terminal 1: API Server
make run-api

# Terminal 2: Scheduler
make run-scheduler
```

### Build for Production
```bash
make build-prod
```

### Testing
```bash
make test
```

## 📊 Job Tracking

All jobs (manual and scheduled) are tracked through the JobManager:

```bash
# Check job status
curl -X GET http://localhost:8080/api/admin/jobs/{jobId} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Job States:**
- `pending` - Job created, waiting to start
- `running` - Job currently executing
- `completed` - Job finished successfully
- `failed` - Job failed with error

## 🔧 Configuration

### Environment Variables
```bash
# Database
DATABASE_URL=postgresql://...

# External API
EXTERNAL_API_URL=https://api.karenai.click
EXTERNAL_API_KEY=your_api_key

# Server
PORT=8080
ENVIRONMENT=development
```

## 📈 Monitoring

### Logs
Both services log to stdout with structured logging:
- API Server: HTTP requests, job creation
- Scheduler: Scheduled runs, job tracking

### Health Checks
```bash
curl http://localhost:8080/health
```

## 🚀 Deployment

### Docker
```bash
# Build image
docker build -t stock-insights .

# Run API server
docker run -p 8080:8080 stock-insights

# Run scheduler (separate container)
docker run stock-insights ./bin/scheduler
```

### Production
```bash
# Build production binaries
make build-prod

# Run with process manager (systemd, supervisor, etc.)
./bin/api     # API server
./bin/scheduler # Scheduler worker
```

## 🔍 Troubleshooting

### Common Issues

1. **Database Connection**
   ```bash
   # Check DATABASE_URL
   echo $DATABASE_URL
   ```

2. **External API**
   ```bash
   # Check API key
   echo $EXTERNAL_API_KEY
   ```

3. **Job Failures**
   ```bash
   # Check job status
   curl -X GET http://localhost:8080/api/admin/jobs/{jobId}
   ```

### Logs
```bash
# API server logs
tail -f logs/api.log

# Scheduler logs  
tail -f logs/scheduler.log
``` 