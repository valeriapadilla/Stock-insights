# Stock Insights MVP - Technical Specifications

## 📋 **Executive Summary**

**Objective**: Stock insights system that supports 10K RPS with robust and scalable architecture.

**Architecture**: CDN → Load Balancer → API Gateway → Services → Database

**Timeline**: 1 week for functional MVP

**Expected RPS**: 1,000-5,000 RPS (Free Tier)
**Maximum RPS**: 10,000 RPS (with optimizations)

---

## 🎯 **1. Functional Requirements**

### **1.1 Data Ingestion**
- ✅ Connect to external API (`api.karenai.click`)
- ✅ Paginated stock ingestion
- ✅ Store data in CockroachDB
- ✅ Daily automatic updates

### **1.2 Public API**
- ✅ List stocks with filters and pagination
- ✅ Get specific stock details
- ✅ Get daily recommendations
- ✅ Basic health check

### **1.3 Admin API**
- ✅ Manual ingestion trigger
- ✅ Recommendation calculation trigger
- ✅ System statistics
- ✅ Detailed health check

### **1.4 UI (Vue3)**
- ✅ Stock list with search
- ✅ Stock details
- ✅ Featured recommendations
- ✅ Responsive design

---

## 🏗️ **2. Technical Architecture**

### **2.1 Architecture Diagram**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   CDN           │    │   Load Balancer │
│   (Vue3)        │───▶│   (Cloudflare)  │───▶│   (Railway)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Monitoring    │    │   API Gateway   │    │   Workers       │
│   (UptimeRobot) │◀───│   (Go)          │◀───│   (Go)         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Cache         │    │   External API  │
                       │   (Memory)      │    │   (karenai.click)│
                       └─────────────────┘    └─────────────────┘
                                │
                       ┌─────────────────┐
                       │   Database      │
                       │   (CockroachDB) │
                       └─────────────────┘
```

### **2.2 Components**

#### **Frontend (Vue3 + TypeScript + Tailwind)**
- **Hosting**: Vercel/Netlify (Free Tier)
- **Framework**: Vue 3 + Composition API
- **State**: Pinia
- **Styling**: Tailwind CSS
- **Build**: Vite

#### **CDN (Cloudflare)**
- **Plan**: Free Tier
- **Functions**: Rate limiting, DDoS protection, caching
- **Rate Limit**: 1000 req/min per IP

#### **Load Balancer (Railway)**
- **Plan**: Free Tier
- **Functions**: Health checks, SSL termination
- **Rate Limit**: 500 req/min per IP

#### **API Gateway (Go)**
- **Framework**: Gin
- **Functions**: Rate limiting, authentication, routing
- **Rate Limit**: 100 req/min per IP

#### **Workers (Go)**
- **Platform**: Heroku Scheduler (Free)
- **Functions**: Ingestion, recommendations
- **Frequency**: Daily

#### **Cache (Memory)**
- **Library**: `patrickmn/go-cache`
- **TTL**: 5 minutes
- **Functions**: Stocks, recommendations

#### **Database (CockroachDB)**
- **Plan**: Cloud Free Tier
- **Functions**: Persistence, complex queries
- **Indexes**: Optimized for performance

#### **Monitoring**
- **Uptime**: UptimeRobot (Free)
- **Metrics**: Grafana Cloud (Free)
- **Logs**: Structured logging

---

## 🔌 **3. API Design**

### **3.1 Public Endpoints (No Auth)**

#### **Health Check**
```bash
GET /api/public/health
# Response:
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### **Stocks**
```bash
GET /api/public/stocks
# Parameters:
# - limit: int (default: 50, max: 100)
# - page: int (default: 1)
# - sort_by: string (company, brokerage, rating_to, time)
# - sort_desc: bool (default: true)
# - q: string (search by ticker/company)
# - brokerage: string (filter by brokerage)
# - rating: string (filter by rating: Buy, Sell, Hold, Neutral)
# - action: string (filter by action: initiated, upgraded, downgraded, etc.)

# Response:
{
  "stocks": [
    {
      "ticker": "AKBA",
      "company": "Akebia Therapeutics",
      "target_from": "$8.00",
      "target_to": "$8.00",
      "rating_from": "Buy",
      "rating_to": "Buy",
      "action": "initiated by",
      "brokerage": "HC Wainwright",
      "time": "2025-06-05T00:30:05.47195313Z",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1000,
    "total_pages": 20
  }
}
```

#### **Stock Details**
```bash
GET /api/public/stocks/{ticker}
# Response:
{
  "ticker": "AKBA",
  "company": "Akebia Therapeutics",
  "target_from": "$8.00",
  "target_to": "$8.00",
  "rating_from": "Buy",
  "rating_to": "Buy",
  "action": "initiated by",
  "brokerage": "HC Wainwright",
  "time": "2025-06-05T00:30:05.47195313Z",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### **Recommendations**
```bash
GET /api/public/recommendations
# Response:
{
  "recommendations": [
    {
      "ticker": "AKBA",
      "company": "Akebia Therapeutics",
      "score": 85.5,
      "rank": 1,
      "explanation": "Strong buy rating from HC Wainwright, target price maintained at $8.00",
      "run_at": "2024-01-15T06:00:00Z"
    }
  ],
  "generated_at": "2024-01-15T06:00:00Z"
}
```

### **3.2 Admin Endpoints (With Auth)**

#### **Detailed Health Check**
```bash
GET /api/admin/health
# Headers: X-API-Key: YOUR_API_KEY
# Response:
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "database_connections": 45,
  "memory_usage": "2.1GB",
  "cache_hit_rate": "95%",
  "active_requests": 1234,
  "uptime": "2h 30m",
  "version": "1.0.0"
}
```

#### **Trigger Ingestion**
```bash
POST /api/admin/ingest
# Headers: X-API-Key: YOUR_API_KEY
# Response:
{
  "status": "started",
  "message": "Ingestion process started",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### **Trigger Recommendations**
```bash
POST /api/admin/recommendations
# Headers: X-API-Key: YOUR_API_KEY
# Response:
{
  "status": "started",
  "message": "Recommendation calculation started",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### **System Stats**
```bash
GET /api/admin/stats
# Headers: X-API-Key: YOUR_API_KEY
# Response:
{
  "total_stocks": 1000,
  "last_ingestion": "2024-01-15T06:00:00Z",
  "last_recommendations": "2024-01-15T06:00:00Z",
  "cache_stats": {
    "hit_rate": "95%",
    "total_items": 1050
  },
  "api_stats": {
    "requests_today": 50000,
    "avg_response_time": "15ms"
  }
}
```

---

## 🗄️ **4. Data Model**

### **4.1 Database Schema**

#### **Stocks Table (Updated according to external API)**
```sql
CREATE TABLE stocks (
    ticker TEXT PRIMARY KEY,
    company TEXT NOT NULL,
    target_from TEXT,
    target_to TEXT,
    rating_from TEXT,
    rating_to TEXT,
    action TEXT,
    brokerage TEXT,
    time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_stocks_brokerage ON stocks(brokerage);      
CREATE INDEX idx_stocks_rating_to ON stocks(rating_to);     
CREATE INDEX idx_stocks_time ON stocks(time DESC);           
CREATE INDEX idx_stocks_company ON stocks(company);          

CREATE INDEX idx_stocks_brokerage_time ON stocks(brokerage, time DESC);
```

#### **Recommendations Table (Updated)**
```sql
CREATE TABLE recommendations (
    id UUID PRIMARY KEY,
    ticker TEXT REFERENCES stocks(ticker),
    score NUMERIC NOT NULL,
    explanation TEXT,
    run_at TIMESTAMPTZ NOT NULL,
    rank INTEGER
);

-- Indexes for recommendations
CREATE INDEX idx_recommendations_run_at ON recommendations(run_at DESC);
CREATE INDEX idx_recommendations_score ON recommendations(score DESC);
CREATE INDEX idx_recommendations_rank ON recommendations(rank);
CREATE INDEX idx_recommendations_ticker ON recommendations(ticker);

-- Composite index for recommendation queries
CREATE INDEX idx_recommendations_run_at_score ON recommendations(run_at DESC, score DESC);
```

### **4.2 Data Types**

#### **Stock (Updated according to external API)**
```go
type Stock struct {
    Ticker      string    `json:"ticker" db:"ticker"`
    Company     string    `json:"company" db:"company"`
    TargetFrom  string    `json:"target_from" db:"target_from"`
    TargetTo    string    `json:"target_to" db:"target_to"`
    RatingFrom  string    `json:"rating_from" db:"rating_from"`
    RatingTo    string    `json:"rating_to" db:"rating_to"`
    Action      string    `json:"action" db:"action"`
    Brokerage   string    `json:"brokerage" db:"brokerage"`
    Time        time.Time `json:"time" db:"time"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// External API Response
type ExternalAPIResponse struct {
    Items []Stock `json:"items"`
}
```

#### **Recommendation**
```go
type Recommendation struct {
    ID          string    `json:"id" db:"id"`
    Ticker      string    `json:"ticker" db:"ticker"`
    Score       float64   `json:"score" db:"score"`
    Explanation string    `json:"explanation" db:"explanation"`
    RunAt       time.Time `json:"run_at" db:"run_at"`
    Rank        int       `json:"rank" db:"rank"`
}
```

### **4.3 API Response Examples**

#### **Stocks List Response**
```json
{
  "stocks": [
    {
      "ticker": "AKBA",
      "company": "Akebia Therapeutics",
      "target_from": "$8.00",
      "target_to": "$8.00",
      "rating_from": "Buy",
      "rating_to": "Buy",
      "action": "initiated by",
      "brokerage": "HC Wainwright",
      "time": "2025-06-05T00:30:05.47195313Z",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "ticker": "CECO",
      "company": "CECO Environmental",
      "target_from": "$33.00",
      "target_to": "$33.00",
      "rating_from": "Neutral",
      "rating_to": "Buy",
      "action": "upgraded by",
      "brokerage": "HC Wainwright",
      "time": "2025-05-01T00:30:06.015697838Z",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1000,
    "total_pages": 20
  }
}
```

#### **Stock Details Response**
```json
{
  "ticker": "AKBA",
  "company": "Akebia Therapeutics",
  "target_from": "$8.00",
  "target_to": "$8.00",
  "rating_from": "Buy",
  "rating_to": "Buy",
  "action": "initiated by",
  "brokerage": "HC Wainwright",
  "time": "2025-06-05T00:30:05.47195313Z",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### **Recommendations Response**
```json
{
  "recommendations": [
    {
      "ticker": "AKBA",
      "company": "Akebia Therapeutics",
      "score": 85.5,
      "rank": 1,
      "explanation": "Strong buy rating from HC Wainwright, target price maintained at $8.00",
      "run_at": "2024-01-15T06:00:00Z"
    },
    {
      "ticker": "CECO",
      "company": "CECO Environmental", 
      "score": 82.3,
      "rank": 2,
      "explanation": "Upgraded to Buy by HC Wainwright, target price $33.00",
      "run_at": "2024-01-15T06:00:00Z"
    }
  ],
  "generated_at": "2024-01-15T06:00:00Z"
}
```

---

## 🚀 **5. Infrastructure**

### **5.1 Free Tier Services**

#### **Database: CockroachDB Cloud**
- **Plan**: Free Tier
- **Storage**: 5GB
- **Connections**: 250
- **Regions**: 1

#### **Backend: Railway**
- **Plan**: Free Tier
- **Instances**: 1
- **Memory**: 512MB
- **Bandwidth**: 100GB/month

#### **Frontend: Vercel**
- **Plan**: Free Tier
- **Bandwidth**: 100GB/month
- **Builds**: Unlimited

#### **CDN: Cloudflare**
- **Plan**: Free Tier
- **Rate Limiting**: Included
- **DDoS Protection**: Included

#### **Monitoring: UptimeRobot**
- **Plan**: Free Tier
- **Checks**: 50
- **Interval**: 5 minutes

### **5.2 Environment Variables**

#### **Backend (.env)**
```bash
# Server Configuration
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info

# Database Configuration
DATABASE_URL=postgres://user:pass@host:26257/stock_insights?sslmode=require

# External API Configuration
EXTERNAL_API_URL=https://api.karenai.click
EXTERNAL_API_KEY=your_api_key_here

# Cache Configuration
CACHE_TTL=5m

# Rate Limiting
RATE_LIMIT=100

# Admin API Key
ADMIN_API_KEY=your_admin_key_here

# Observability
METRICS_PORT=9090
ENABLE_TRACING=false
```

---

## 🔒 **6. Security**

### **6.1 Rate Limiting**
- **CDN**: 1000 req/min per IP
- **Load Balancer**: 500 req/min per IP
- **API Gateway**: 100 req/min per IP

### **6.2 Authentication**
- **Public**: No auth (with rate limiting)
- **Admin**: API Key in header `X-API-Key`

### **6.3 Security Headers**
```go
// Automatic headers
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "DENY")
c.Header("X-XSS-Protection", "1; mode=block")
c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
```

### **6.4 CORS**
```go
// CORS configuration
config := cors.DefaultConfig()
config.AllowOrigins = []string{"https://your-frontend.vercel.app"}
config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
```

---

## 📊 **7. Monitoring & Observability**

### **7.1 Health Checks**
- **Public**: `/api/public/health` (basic)
- **Admin**: `/api/admin/health` (detailed)

### **7.2 Metrics**
- **Requests/sec**: Prometheus
- **Response time**: P95 < 300ms
- **Error rate**: < 1%
- **Cache hit rate**: > 90%

### **7.3 Logging**
```go
// Structured logging
log.WithFields(log.Fields{
    "method":     c.Request.Method,
    "path":       c.Request.URL.Path,
    "ip":         c.ClientIP(),
    "user_agent": c.Request.UserAgent(),
    "duration":   time.Since(start),
}).Info("Request completed")
```

---

## 🧪 **8. Testing Strategy**

### **8.1 Unit Tests**
- **Coverage**: > 80%
- **Services**: Ingestion, recommendations
- **Cache**: Memory cache operations
- **Repository**: Database operations

### **8.2 Integration Tests**
- **API endpoints**: All endpoints
- **Database**: CRUD operations
- **External API**: Mock responses

### **8.3 Load Tests**
- **Tool**: k6
- **Target**: 1000 RPS (Free Tier), 5000 RPS (maximum)
- **Duration**: 5 minutes
- **Success criteria**: P95 < 300ms

---

🗄️ Fase 2: Database & Repository (2 horas)
2.1 Configurar CockroachDB
✅ Crear: internal/database/connection.go
✅ Test: Conexión a CockroachDB
✅ Test: Ping a la base de datos
2.2 Migraciones
✅ Crear: internal/database/migrations/
✅ SQL: Tabla stocks (según specs)
✅ SQL: Tabla recommendations (según specs)
✅ Test: Migraciones se ejecutan correctamente
2.3 Repository Pattern
✅ Crear: internal/repository/stock_repository.go
✅ Test: CRUD básico para stocks
✅ Crear: internal/repository/recommendation_repository.go
✅ Test: CRUD básico para recommendations
🔌 Fase 3: External API & Data Worker (3 horas)
3.1 External API Client
✅ Crear: internal/client/external_api.go
✅ Test: Conectar a api.karenai.click
✅ Test: Manejar paginación
✅ Test: Manejar errores de red
3.2 Data Worker (Cada 24 horas)
✅ Crear: cmd/worker/data/main.go
✅ Crear: internal/worker/data_worker.go
✅ Test: Ingestion completa de datos
✅ Test: Manejo de duplicados (upsert)
✅ Test: Logging de progreso
✅ Test: Ejecución manual y automática
3.3 Ingestion Endpoint (Manual trigger)
✅ Crear: internal/handler/ingestion.go
✅ Test: POST /api/admin/ingest (con JWT)
✅ Test: Respuesta con status y progreso
🌐 Fase 4: API Gateway & Middleware (1 hora)
4.1 API Gateway Pattern
✅ Crear: internal/gateway/api_gateway.go
✅ Test: Request/response middleware
✅ Test: CORS middleware
✅ Test: Security headers middleware
✅ Test: Logging middleware
✅ Test: Metrics middleware
4.2 Cache Implementation
✅ Crear: internal/cache/memory_cache.go
✅ Test: Cache hit/miss metrics
✅ Test: Cache invalidation strategy
✅ Test: Cache warming for hot data
✅ Test: Cache middleware
✅ Test: TTL configuration (stocks: 5min, details: 10min, recommendations: 1h)
✅ Test: Cache hit rate monitoring
🔌 Fase 5: API Server - Stock Endpoints (2 horas)
5.1 Stock Service
✅ Crear: internal/service/stock_service.go
✅ Test: Obtener stocks con filtros
✅ Test: Paginación
✅ Test: Búsqueda por ticker/company
5.2 Stock Handlers
✅ Crear: internal/handler/stock.go
✅ Test: GET /api/public/stocks
✅ Test: GET /api/public/stocks/{ticker}
✅ Test: Query parameters (sort, filter, page)
5.3 Cache Integration
✅ Integrar: Cache en stock service
✅ Test: Cache de stocks
✅ Test: TTL y invalidación
🤖 Fase 6: Recommendation Worker & Algorithm (3 horas)
6.1 Recommendation Worker (Cada 24 horas)
✅ Crear: cmd/worker/recommendation/main.go
✅ Crear: internal/worker/recommendation_worker.go
✅ Test: Cálculo automático diario
✅ Test: Algoritmo de scoring básico
✅ Test: Ranking de stocks
✅ Test: Generación de explicaciones
6.2 Recommendation Service
✅ Crear: internal/service/recommendation_service.go
✅ Test: Algoritmo de scoring
✅ Test: Ranking y explicaciones
✅ Test: Cache de recommendations
6.3 Recommendation Handlers
✅ Crear: internal/handler/recommendation.go
✅ Test: GET /api/public/recommendations
✅ Test: Respuesta con top stocks
6.4 Recommendation Endpoint (Manual trigger)
✅ Test: POST /api/admin/recommendations (con JWT)
✅ Test: Trigger manual de cálculo
