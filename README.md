# 📈 StockInsights

A modern stock analysis and recommendation platform built with Go backend and Vue.js frontend.

## 🚀 Features

- **Real-time Stock Data**: Track and analyze stock market data
- **AI-Powered Recommendations**: Get daily stock recommendations based on advanced algorithms
- **Advanced Filtering**: Filter stocks by rating, sort by ticker or change percentage
- **Modern UI**: Dark theme with responsive design
- **Real-time Updates**: Live data with last update indicators

## 🛠️ Tech Stack

### Backend
- **Go 1.23** - High-performance server language
- **Gin** - Fast HTTP web framework
- **CockroachDB** - Distributed SQL database
- **JWT** - Authentication and authorization
- **Logrus** - Structured logging
- **Testify** - Testing framework

### Frontend
- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe JavaScript
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Tailwind CSS** - Utility-first CSS framework
- **Axios** - HTTP client
- **Vite** - Fast build tool

## 🌐 Live Demo

- **Frontend**: [https://stock-insights-sigma.vercel.app](https://stock-insights-sigma.vercel.app)
- **Backend API**: [https://stock-insights-production-3f39.up.railway.app](https://stock-insights-production-3f39.up.railway.app)

## 🚀 Quick Start

### Prerequisites

- **Go 1.23+** - [Download](https://golang.org/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **CockroachDB** - [Download](https://www.cockroachlabs.com/docs/stable/install-cockroachdb.html)

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Stock-insights
   ```

2. **Configure environment variables**
   ```bash
   cd backend
   cp env.example .env
   # Edit .env with your configuration
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run database migrations**
   ```bash
   go run cmd/migrate/main.go
   ```

5. **Start the backend server**
   ```bash
   go run cmd/api/main.go
   ```

The backend will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Configure environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your API configuration
   ```

3. **Install dependencies**
   ```bash
   npm install
   ```

4. **Start development server**
   ```bash
   npm run dev
   ```

The frontend will be available at `http://localhost:5173`

## 🚀 Deployment

### Backend Deployment
```bash
cd backend
go build -o bin/stock-insights cmd/api/main.go
```

#### Frontend Deployment
```bash
cd frontend
npm run build
```
