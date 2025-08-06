# StockInsights Frontend

Modern Vue.js frontend for the StockInsights platform featuring a responsive dark theme and real-time stock data visualization.

## 🎨 Features

- **Dark Theme UI**: Modern, professional interface
- **Real-time Data**: Live stock updates with timestamps
- **Advanced Filtering**: Filter by rating and sort options
- **Responsive Design**: Works on desktop and mobile
- **Stock Details Modal**: Detailed view for each stock
- **Recommendations**: AI-powered stock recommendations
- **Infinite Scroll**: Load more stocks seamlessly

## 🛠️ Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe JavaScript
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Tailwind CSS** - Utility-first CSS framework
- **Axios** - HTTP client
- **Vite** - Fast build tool

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- npm or yarn
- Backend API running

### Installation

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Copy environment file**
   ```bash
   cp env.example .env
   ```

3. **Configure environment variables**
   ```bash
   # Edit .env with your API configuration
   VITE_API_BASE_URL=http://localhost:8080/api/v1
   ```

4. **Install dependencies**
   ```bash
   npm install
   ```

5. **Start development server**
   ```bash
   npm run dev
   ```

The application will be available at `http://localhost:5173`

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/         # Vue components
│   │   ├── common/        # Shared components
│   │   ├── dashboard/     # Dashboard components
│   │   ├── stocks/        # Stock-related components
│   │   └── recommendations/ # Recommendation components
│   ├── stores/            # Pinia state management
│   │   ├── stocks.ts      # Stock store
│   │   └── recommendations.ts # Recommendations store
│   ├── services/          # API services
│   │   ├── api.ts         # Axios configuration
│   │   └── stocks.ts      # Stock API calls
│   ├── types/             # TypeScript type definitions
│   │   └── api.ts         # API response types
│   ├── utils/             # Utility functions
│   │   ├── sort.ts        # Sorting utilities
│   │   └── stock.ts       # Stock utilities
│   ├── views/             # Page components
│   │   └── DashboardView.vue # Main dashboard
│   ├── router/            # Vue Router configuration
│   │   └── index.ts       # Route definitions
│   ├── App.vue            # Root component
│   └── main.ts            # Application entry point
├── public/                # Static assets
└── package.json           # Dependencies and scripts
```

## 🎯 Components Overview

### Core Components

- **Header**: Navigation and status indicator
- **SummaryCards**: Dashboard metrics display
- **FilterBar**: Stock filtering and search
- **StockCard**: Individual stock display
- **RecommendationCard**: Recommendation display
- **Modal**: Reusable modal component
- **StockDetails**: Detailed stock information

### State Management

- **Stocks Store**: Manages stock data and filters
- **Recommendations Store**: Manages recommendation data

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | Backend API URL | `http://localhost:8080/api/v1` |
| `VITE_APP_ENV` | Environment | `development` |

### API Integration

The frontend communicates with the backend through:

- **Stocks API**: `/api/v1/public/stocks`
- **Search API**: `/api/v1/public/stocks/search`
- **Recommendations API**: `/api/v1/public/recommendations`

## 🎨 UI/UX Features

### Dark Theme
- Professional dark color scheme
- High contrast for readability
- Consistent color palette

### Responsive Design
- Mobile-first approach
- Flexible grid layouts
- Touch-friendly interactions

### Loading States
- Skeleton loading for cards
- Spinner animations
- Error handling with retry options

### Animations
- Smooth transitions
- Modal animations
- Hover effects

## 🧪 Testing

### Run tests
```bash
npm run test
```

### Run with coverage
```bash
npm run test:coverage
```

### Run e2e tests
```bash
npm run test:e2e
```

## 🚀 Build & Deployment

### Development
```bash
npm run dev
```

### Production build
```bash
npm run build
```

### Preview production build
```bash
npm run preview
```

### Docker deployment
```bash
# Build image
docker build -t stock-insights-frontend .

# Run container
docker run -p 80:80 stock-insights-frontend
```

## 🔍 Development Tools

### Vue DevTools
Install Vue DevTools browser extension for debugging.

### TypeScript
Full TypeScript support with strict type checking.

### ESLint & Prettier
Code formatting and linting configured.

## 🐛 Troubleshooting

### Common Issues

1. **API connection failed**
   - Check `VITE_API_BASE_URL` in `.env`
   - Ensure backend is running
   - Check CORS configuration

2. **Build errors**
   - Clear node_modules: `rm -rf node_modules && npm install`
   - Check TypeScript errors
   - Verify all dependencies

3. **Hot reload not working**
   - Check file watchers
   - Restart dev server
   - Clear browser cache

### Debug Mode

Enable debug logging:
```bash
# Add to .env
VITE_DEBUG=true
```

