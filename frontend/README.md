# StockInsights Frontend

A modern Vue 3 application for displaying stock insights and recommendations with a beautiful dark theme UI.

## Quick Start

### Prerequisites
- Node.js 18+
- Backend API running 

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The application will be available at `http://localhost:5173`

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type safety and better developer experience
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Tailwind CSS** - Utility-first CSS framework
- **Axios** - HTTP client for API calls
- **Vite** - Fast build tool

## Project Structure

```
src/
├── components/
│   ├── common/          # Shared components (Header, SearchBar, etc.)
│   ├── dashboard/       # Dashboard components (SummaryCards)
│   ├── stocks/          # Stock-related components
│   └── recommendations/ # Recommendation components
├── stores/              # Pinia stores (stocks, recommendations)
├── services/            # API services
├── types/               # TypeScript interfaces
├── views/               # Page components
├── router/              # Vue Router configuration
└── utils/               # Utility functions and constants
```

## Features

### Dashboard
- **Summary Cards** - Key metrics (Total Stocks, Upgrades, Downgrades, Recommendations)
- **Live Data Indicator** - Real-time data status
- **Navigation Tabs** - Switch between Stocks and Recommendations views

### Stocks View
- **Search & Filter** - Search by ticker/company, filter by rating/action
- **Stock Cards** - Display stock information with price changes
- **Pagination** - Load more stocks as needed

### Recommendations View
- **AI Recommendations** - Top stock recommendations with scores
- **Ranking System** - Ranked by AI algorithm score
- **Detailed Explanations** - AI-generated reasoning for each recommendation

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the frontend directory:

```env
VITE_API_URL=http://localhost:8080/api/v1/public
```

### API Integration

The frontend connects to the StockInsights backend API:

- **Stocks API** - `/api/v1/public/stocks`
- **Recommendations API** - `/api/v1/public/recommendations`
- **Search API** - `/api/v1/public/stocks/search`

## Key Components

### Header Component
- StockInsights logo with green arrow icon
- Live data indicator with pulsing animation

### SummaryCards Component
- Displays key metrics in card format
- Loading states and error handling
- Responsive grid layout

### FilterBar Component
- Search functionality with debouncing
- Dropdown filters for rating, action, and price range
- Active filter indicators with clear options

### StockCard Component
- Displays stock information in card format
- Color-coded rating badges
- Price change indicators with arrows
- Hover effects and transitions

### RecommendationCard Component
- AI score display with gradient background
- Ranking system (#1, #2, etc.)
- Detailed explanations
- View details functionality

## 🚀 Development

```bash
# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Design System

The application uses a dark theme with:

- **Primary Colors**: Green (#10B981) for positive actions and highlights
- **Background**: Dark gray (#111827) for main background
- **Cards**: Medium gray (#1F2937) for card backgrounds
- **Text**: White for primary text, gray for secondary text
- **Accents**: Red for negative changes, yellow for neutral states

## Responsive Design

The application is fully responsive and works on:
- Desktop (1200px+)
- Tablet (768px - 1199px)
- Mobile (320px - 767px)

## Backend Integration

This frontend is designed to work with the StockInsights backend API. Make sure the backend is running and accessible at the configured API URL.
