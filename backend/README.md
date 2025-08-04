# Stock Insights Backend

A high-performance stock insights system that provides real-time stock data and recommendations. Built with Go, featuring a robust architecture with automatic data ingestion, recommendation algorithms, and comprehensive testing.

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- CockroachDB
- Environment variables configured

### Installation
```bash
# Install dependencies
make deps

# Build the application
make build

# Run tests
make test
```

## 📦 Architecture Overview

```
AGRGAR FOTO
```

## 🎯 Core Features

### **Data Ingestion**
- ✅ Automatic daily stock data ingestion from external API
- ✅ Manual ingestion trigger via admin API
- ✅ Efficient incremental updates (only new data)
- ✅ Job tracking and monitoring

### **Stock Recommendations**
- ✅ Recommendation algorithm (0-100 scoring)
- ✅ Daily automatic recommendation calculation
- ✅ Manual recommendation trigger via admin API
- ✅ Data freshness validation

### **Public API**
- ✅ Stock listing with advanced filtering and pagination
- ✅ Individual stock details
- ✅ Stock search functionality
- ✅ Daily recommendations with scoring

### **Admin API**
- ✅ Manual data ingestion trigger
- ✅ Manual recommendation calculation
- ✅ Job status tracking
- ✅ System health monitoring

## 🔧 Components

### 1. API Server
The main API server handling HTTP requests with authentication and rate limiting.

```bash
# Run API server
make run-api

# Or directly
go run cmd/api/main.go
```

### 2. Data Ingestion Worker
Automatically runs stock data ingestion every 24 hours.

```bash
# Run ingestion worker
make run-scheduler

# Or directly
go run cmd/worker/scheduler/main.go
```

### 3. Recommendation Worker
Calculates daily stock recommendations based on scoring algorithm.

```bash
# Run recommendation worker
make run-recommendations

# Or directly
go run cmd/worker/recommendations/main.go
```

## 📡 API Endpoints

### Public Endpoints

#### **Stocks**
```bash
# List stocks with pagination and filters
GET /api/v1/public/stocks?limit=10&offset=0&sort=time&order=desc

# Get specific stock
GET /api/v1/public/stocks/{ticker}

# Search stocks with filters
GET /api/v1/public/stocks/search?ticket=AAPL&date_from=2025-01-01
```

#### **Recommendations**
```bash
# Get daily recommendations
GET /api/v1/public/recommendations?limit=10
```

#### **Health Check**
```bash
# System health
GET /health
GET /api/v1/public/health
```

### Admin Endpoints (Require Authentication)

#### **Data Ingestion**
```bash
# Trigger manual ingestion
POST /api/v1/admin/ingest/stocks
Authorization: Bearer <admin_token>

# Check job status
GET /api/v1/admin/jobs/{jobId}
Authorization: Bearer <admin_token>
```

#### **Recommendations**
```bash
# Calculate recommendations manually
POST /api/v1/admin/recommendations/calculate
Authorization: Bearer <admin_token>
```

## 🧠 Recommendation Algorithm

The system uses a sophisticated scoring algorithm (0-100 points):

### **Scoring Components**
- **Action Score (0-40)**: Based on analyst actions (target raised/lowered)
- **Rating Score (0-25)**: Based on analyst ratings (Buy, Overweight, etc.)
- **Target Change Score (0-20)**: Based on percentage change in target price
- **Freshness Score (0-15)**: Based on data recency

### **Filtering Process**
1. **Freshness Filter**: Only stocks from last 7 days
2. **Positive Action Filter**: Only stocks with positive analyst actions >80
3. **Positive Rating Filter**: Only stocks with positive ratings
4. **Significant Change Filter**: Only stocks with significant target changes

## 🛠️ Development

### Project Structure
```
backend/
├── cmd/                    # Application entry points
│   ├── api/               # API server
│   ├── worker/            # Workers (ingestion, recommendations)
│   ├── migrate/           # Database migrations
│   └── setup-auth/        # Setup authentication for admin
├── internal/              # Internal packages
│   ├── app/               # Application setup
│   ├── bootstrap/         # Bootstrap configuration
│   ├── client/            # External API client
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and migrations
│   ├── dto/               # Data Transfer Objects
│   ├── errors/            # Custom error types
│   ├── handler/           # HTTP handlers
│   ├── job/               # Job management
│   ├── middleware/        # HTTP middleware
│   ├── model/             # Data models
│   ├── repository/        # Data access layer
│   ├── server/            # Server setup
│   ├── service/           # Business logic
│   ├── validator/         # Input validation
│   └── worker/            # Background workers
├── scripts/               # Utility scripts
├── docs/                  # Documentation
├── Dockerfile             # Docker configuration
└── Makefile               # Build and development commands
```
## 🔧 Configuration

### Environment Variables
```bash
# Database
DATABASE_URL=postgresql://user:password@host:port/database

# External API
EXTERNAL_API_URL=https://api.karenai.click
EXTERNAL_API_KEY=your_api_key

# Server
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

# Admin Authentication
ADMIN_API_KEY=your_admin_key

# Rate Limiting
RATE_LIMIT=100

# Caching
CACHE_TTL=5m
```

## 📊 Job Tracking

All jobs (manual and scheduled) are tracked through the JobManager:

### Job States
- `pending` - Job created, waiting to start
- `running` - Job currently executing
- `completed` - Job finished successfully
- `failed` - Job failed with error

### Job Monitoring
```bash
# Check job status
curl -X GET http://localhost:8080/api/v1/admin/jobs/{jobId} \
  -H "Authorization: Bearer YOUR_TOKEN"
```
## 🚀 Deployment

### Docker Deployment
```bash
# Build image
docker build -t stock-insights .

# Run API server
docker run -p 8080:8080 stock-insights

# Run workers (separate containers)
docker run stock-insights ./bin/scheduler
docker run stock-insights ./bin/recommendations
```

### GitHub Actions Scheduling
The system uses GitHub Actions for scheduled workers:

- **Ingestion Worker**: Runs daily at 4:00 AM Colombia time
- **Recommendation Worker**: Runs daily at 4:30 AM Colombia time

### Performance Optimization
- Database queries are optimized with proper indexing
- Incremental data ingestion (only new data)
- Efficient recommendation calculation with filtering
- Rate limiting to prevent abuse

## 🧪 Testing Coverage

The project maintains comprehensive test coverage:

- **Service Layer**: 68.8% coverage
- **Handler Layer**: 81.7% coverage
- **Validator Layer**: 96.6% coverage
- **Client Layer**: 91.4% coverage
- **Config Layer**: 80.0% coverage