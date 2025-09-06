# Phase 5: MVP Polish & Deployment

## Overview
Final phase to polish the MVP, add essential production features, and prepare for deployment. Focus on reliability, user experience, and operational readiness.

## Goals
- Add essential production features (logging, monitoring, error handling)
- Improve user experience with better UI/UX
- Add deployment configuration and documentation
- Implement basic analytics and admin features
- Ensure system is ready for real users
- Create comprehensive documentation

## Tasks

### 5.1 Production Logging & Monitoring
**File**: `internal/infra/logging/logger.go` (new)

```go
package logging

import (
    "context"
    "encoding/json"
    "log/slog"
    "os"
    "time"
)

type Logger struct {
    *slog.Logger
}

type LogEntry struct {
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Timestamp time.Time              `json:"timestamp"`
    UserID    *int64                 `json:"user_id,omitempty"`
    QuestID   *string                `json:"quest_id,omitempty"`
    Operation string                 `json:"operation,omitempty"`
    Duration  *time.Duration         `json:"duration,omitempty"`
    Error     string                 `json:"error,omitempty"`
    Extra     map[string]interface{} `json:"extra,omitempty"`
}

func NewLogger() *Logger {
    opts := &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }
    
    if os.Getenv("LOG_FORMAT") == "json" {
        handler := slog.NewJSONHandler(os.Stdout, opts)
        return &Logger{slog.New(handler)}
    }
    
    handler := slog.NewTextHandler(os.Stdout, opts)
    return &Logger{slog.New(handler)}
}

func (l *Logger) LogQuestCompletion(ctx context.Context, userID int64, questID string, pointsAwarded string, duration time.Duration) {
    l.Info("quest_completed",
        "user_id", userID,
        "quest_id", questID,
        "points_awarded", pointsAwarded,
        "duration_ms", duration.Milliseconds(),
        "operation", "quest_completion",
    )
}

func (l *Logger) LogUserRegistration(ctx context.Context, userID int64, source string) {
    l.Info("user_registered",
        "user_id", userID,
        "source", source,
        "operation", "user_registration",
    )
}

func (l *Logger) LogError(ctx context.Context, operation string, err error, extra map[string]interface{}) {
    args := []interface{}{
        "operation", operation,
        "error", err.Error(),
    }
    
    for k, v := range extra {
        args = append(args, k, v)
    }
    
    l.Error("operation_failed", args...)
}

func (l *Logger) LogAPIRequest(ctx context.Context, method, path string, userID *int64, duration time.Duration, statusCode int) {
    args := []interface{}{
        "method", method,
        "path", path,
        "duration_ms", duration.Milliseconds(),
        "status_code", statusCode,
        "operation", "api_request",
    }
    
    if userID != nil {
        args = append(args, "user_id", *userID)
    }
    
    if statusCode >= 400 {
        l.Warn("api_request_error", args...)
    } else {
        l.Info("api_request", args...)
    }
}
```

**Testing**:
- Test structured logging output
- Test different log levels
- Test context propagation
- Test log formatting (JSON vs text)

### 5.2 Health Check & Metrics Endpoints
**File**: `internal/infra/http/health_handler.go` (new)

```go
package http

import (
    "context"
    "database/sql"
    "encoding/json"
    "net/http"
    "time"
)

type HealthHandler struct {
    db *sql.DB
}

type HealthResponse struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Version   string            `json:"version"`
    Checks    map[string]Check  `json:"checks"`
}

type Check struct {
    Status  string        `json:"status"`
    Message string        `json:"message,omitempty"`
    Latency time.Duration `json:"latency_ms"`
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
    return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    response := HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   getVersion(),
        Checks:    make(map[string]Check),
    }
    
    // Database check
    dbStart := time.Now()
    err := h.db.PingContext(ctx)
    dbLatency := time.Since(dbStart)
    
    if err != nil {
        response.Status = "unhealthy"
        response.Checks["database"] = Check{
            Status:  "unhealthy",
            Message: err.Error(),
            Latency: dbLatency,
        }
    } else {
        response.Checks["database"] = Check{
            Status:  "healthy",
            Latency: dbLatency,
        }
    }
    
    // Set response status
    statusCode := http.StatusOK
    if response.Status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
    // Simple readiness check
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ready",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

func getVersion() string {
    version := os.Getenv("APP_VERSION")
    if version == "" {
        return "dev"
    }
    return version
}
```

**Testing**:
- Test health endpoint returns correct status
- Test database connectivity check
- Test readiness endpoint
- Test timeout handling

### 5.3 Enhanced Error Handling & User Feedback
**File**: `internal/infra/http/middleware.go` (new)

```go
package http

import (
    "context"
    "net/http"
    "time"
    
    "github.com/supercakecrumb/adhd-game-bot/internal/infra/logging"
)

type Middleware struct {
    logger *logging.Logger
}

func NewMiddleware(logger *logging.Logger) *Middleware {
    return &Middleware{logger: logger}
}

func (m *Middleware) RequestLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start)
        
        // Extract user ID from context if available
        var userID *int64
        if uid := r.Context().Value("user_id"); uid != nil {
            if id, ok := uid.(int64); ok {
                userID = &id
            }
        }
        
        m.logger.LogAPIRequest(r.Context(), r.Method, r.URL.Path, userID, duration, wrapped.statusCode)
    })
}

func (m *Middleware) ErrorRecovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                m.logger.LogError(r.Context(), "panic_recovery", fmt.Errorf("%v", err), map[string]interface{}{
                    "method": r.Method,
                    "path":   r.URL.Path,
                })
                
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // Allow specific origins in production
        allowedOrigins := []string{
            "http://localhost:3000",
            "https://yourdomain.com",
        }
        
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

**Testing**:
- Test request logging captures all details
- Test panic recovery works
- Test CORS headers are set correctly
- Test status code capture

### 5.4 Improved Frontend Error Handling
**File**: `frontend/src/components/common/ErrorBoundary.tsx` (new)

```typescript
import React, { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Uncaught error:', error, errorInfo);
    
    // In production, send to error tracking service
    if (process.env.NODE_ENV === 'production') {
      // Send to Sentry, LogRocket, etc.
    }
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen bg-slate-900 flex items-center justify-center">
          <div className="bg-slate-800 p-8 rounded-lg max-w-md w-full mx-4">
            <div className="text-center">
              <div className="text-red-500 text-6xl mb-4">⚠️</div>
              <h1 className="text-2xl font-bold text-slate-100 mb-4">
                Something went wrong
              </h1>
              <p className="text-slate-400 mb-6">
                We're sorry, but something unexpected happened. Please try refreshing the page.
              </p>
              <div className="space-y-2">
                <button
                  onClick={() => window.location.reload()}
                  className="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded"
                >
                  Refresh Page
                </button>
                <button
                  onClick={() => window.location.href = '/'}
                  className="w-full bg-slate-600 hover:bg-slate-700 text-white py-2 px-4 rounded"
                >
                  Go Home
                </button>
              </div>
              {process.env.NODE_ENV === 'development' && this.state.error && (
                <details className="mt-4 text-left">
                  <summary className="text-slate-400 cursor-pointer">Error Details</summary>
                  <pre className="text-xs text-red-400 mt-2 overflow-auto">
                    {this.state.error.stack}
                  </pre>
                </details>
              )}
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
```

**File**: `frontend/src/hooks/useErrorHandler.ts` (new)

```typescript
import { useState, useCallback } from 'react';

interface ErrorState {
  error: string | null;
  isLoading: boolean;
}

export const useErrorHandler = () => {
  const [state, setState] = useState<ErrorState>({
    error: null,
    isLoading: false,
  });

  const handleAsync = useCallback(async <T>(
    asyncFn: () => Promise<T>,
    errorMessage?: string
  ): Promise<T | null> => {
    setState({ error: null, isLoading: true });
    
    try {
      const result = await asyncFn();
      setState({ error: null, isLoading: false });
      return result;
    } catch (error) {
      const message = errorMessage || 
        (error instanceof Error ? error.message : 'An unexpected error occurred');
      
      setState({ error: message, isLoading: false });
      return null;
    }
  }, []);

  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }));
  }, []);

  return {
    error: state.error,
    isLoading: state.isLoading,
    handleAsync,
    clearError,
  };
};
```

**Testing**:
- Test error boundary catches and displays errors
- Test error handler hook manages async errors
- Test error reporting in production
- Test user-friendly error messages

### 5.5 Basic Analytics & Admin Dashboard
**File**: `internal/usecase/analytics_service.go` (new)

```go
package usecase

import (
    "context"
    "time"
    
    "github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
    "github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type AnalyticsService struct {
    questCompletionRepo ports.QuestCompletionRepository
    userRepo           ports.UserRepository
    questRepo          ports.QuestRepository
    dungeonRepo        ports.DungeonRepository
}

type SystemStats struct {
    TotalUsers          int                     `json:"total_users"`
    ActiveUsers         int                     `json:"active_users_7d"`
    TotalQuests         int                     `json:"total_quests"`
    TotalCompletions    int                     `json:"total_completions"`
    TotalPointsAwarded  valueobject.Decimal     `json:"total_points_awarded"`
    AvgCompletionsPerUser float64               `json:"avg_completions_per_user"`
    TopQuests           []QuestStats            `json:"top_quests"`
    RecentActivity      []ActivityItem          `json:"recent_activity"`
}

type QuestStats struct {
    QuestID       string              `json:"quest_id"`
    Title         string              `json:"title"`
    Completions   int                 `json:"completions"`
    TotalPoints   valueobject.Decimal `json:"total_points"`
    AvgPoints     valueobject.Decimal `json:"avg_points"`
}

type ActivityItem struct {
    Type        string    `json:"type"`
    Description string    `json:"description"`
    Timestamp   time.Time `json:"timestamp"`
    UserID      int64     `json:"user_id"`
    Username    string    `json:"username"`
}

func NewAnalyticsService(
    questCompletionRepo ports.QuestCompletionRepository,
    userRepo ports.UserRepository,
    questRepo ports.QuestRepository,
    dungeonRepo ports.DungeonRepository,
) *AnalyticsService {
    return &AnalyticsService{
        questCompletionRepo: questCompletionRepo,
        userRepo:           userRepo,
        questRepo:          questRepo,
        dungeonRepo:        dungeonRepo,
    }
}

func (s *AnalyticsService) GetSystemStats(ctx context.Context) (*SystemStats, error) {
    stats := &SystemStats{}
    
    // Get total users
    totalUsers, err := s.userRepo.GetTotalCount(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalUsers = totalUsers
    
    // Get active users (completed quest in last 7 days)
    activeUsers, err := s.questCompletionRepo.GetActiveUsersCount(ctx, 7)
    if err != nil {
        return nil, err
    }
    stats.ActiveUsers = activeUsers
    
    // Get total quests
    totalQuests, err := s.questRepo.GetTotalCount(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalQuests = totalQuests
    
    // Get completion stats
    completionStats, err := s.questCompletionRepo.GetGlobalStats(ctx)
    if err != nil {
        return nil, err
    }
    stats.TotalCompletions = completionStats.TotalCompletions
    stats.TotalPointsAwarded = completionStats.TotalPoints
    
    if totalUsers > 0 {
        stats.AvgCompletionsPerUser = float64(stats.TotalCompletions) / float64(totalUsers)
    }
    
    // Get top quests
    topQuests, err := s.questCompletionRepo.GetTopQuests(ctx, 10)
    if err != nil {
        return nil, err
    }
    stats.TopQuests = topQuests
    
    // Get recent activity
    recentActivity, err := s.questCompletionRepo.GetRecentActivity(ctx, 20)
    if err != nil {
        return nil, err
    }
    stats.RecentActivity = recentActivity
    
    return stats, nil
}

func (s *AnalyticsService) GetDungeonStats(ctx context.Context, dungeonID string) (*DungeonStats, error) {
    // Implementation for dungeon-specific stats
    return nil, nil
}
```

**File**: `frontend/src/pages/Admin.tsx` (modify existing)

```typescript
import React, { useState, useEffect } from 'react';
import { Users, Target, Award, TrendingUp } from 'lucide-react';
import { apiService } from '../services/api';

interface SystemStats {
  total_users: number;
  active_users_7d: number;
  total_quests: number;
  total_completions: number;
  total_points_awarded: string;
  avg_completions_per_user: number;
  top_quests: QuestStats[];
  recent_activity: ActivityItem[];
}

interface QuestStats {
  quest_id: string;
  title: string;
  completions: number;
  total_points: string;
  avg_points: string;
}

interface ActivityItem {
  type: string;
  description: string;
  timestamp: string;
  user_id: number;
  username: string;
}

const Admin: React.FC<{ user: any }> = ({ user }) => {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const systemStats = await apiService.getSystemStats();
      setStats(systemStats);
    } catch (error) {
      console.error('Failed to load stats:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-slate-400">Loading analytics...</div>
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-400">Failed to load analytics</div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-slate-100 mb-8">Admin Dashboard</h1>
      
      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-slate-800 p-6 rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-slate-400 text-sm">Total Users</p>
              <p className="text-2xl font-bold text-slate-100">{stats.total_users}</p>
            </div>
            <Users className="w-8 h-8 text-blue-500" />
          </div>
        </div>
        
        <div className="bg-slate-800 p-6 rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-slate-400 text-sm">Active Users (7d)</p>
              <p className="text-2xl font-bold text-slate-100">{stats.active_users_7d}</p>
            </div>
            <TrendingUp className="w-8 h-8 text-green-500" />
          </div>
        </div>
        
        <div className="bg-slate-800 p-6 rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-slate-400 text-sm">Total Quests</p>
              <p className="text-2xl font-bold text-slate-100">{stats.total_quests}</p>
            </div>
            <Target className="w-8 h-8 text-purple-500" />
          </div>
        </div>
        
        <div className="bg-slate-800 p-6 rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-slate-400 text-sm">Total Completions</p>
              <p className="text-2xl font-bold text-slate-100">{stats.total_completions}</p>
            </div>
            <Award className="w-8 h-8 text-yellow-500" />
          </div>
        </div>
      </div>
      
      {/* Top Quests */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-slate-800 p-6 rounded-lg">
          <h2 className="text-xl font-bold text-slate-100 mb-4">Top Quests</h2>
          <div className="space-y-3">
            {stats.top_quests.slice(0, 5).map((quest, index) => (
              <div key={quest.quest_id} className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="w-6 h-6 bg-blue-600 rounded-full flex items-center justify-center text-xs font-bold">
                    {index + 1}
                  </div>
                  <div>
                    <p className="text-slate-100 font-medium">{quest.title}</p>
                    <p className="text-slate-400 text-sm">{quest.completions} completions</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-slate-100 font-medium">{quest.total_points} pts</p>
                  <p className="text-slate-400 text-sm">avg: {quest.avg_points}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
        
        {/* Recent Activity */}
        <div className="bg-slate-800 p-6 rounded-lg">
          <h2 className="text-xl font-bold text-slate-100 mb-4">Recent Activity</h2>
          <div className="space-y-3">
            {stats.recent_activity.slice(0, 10).map((activity, index) => (
              <div key={index} className="flex items-start space-x-3">
                <div className="w-2 h-2 bg-green-500 rounded-full mt-2"></div>
                <div className="flex-1">
                  <p className="text-slate-100 text-sm">{activity.description}</p>
                  <p className="text-slate-400 text-xs">
                    {new Date(activity.timestamp).toLocaleString()}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Admin;
```

**Testing**:
- Test analytics data collection
- Test admin dashboard displays correctly
- Test real-time stats updates
- Test performance with large datasets

### 5.6 Deployment Configuration
**File**: `docker-compose.prod.yml` (new)

```yaml
version: '3.8'

services:
  db:
    image: postgres:13
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 30s
      timeout: 10s
      retries: 3

  api:
    build: 
      context: .
      dockerfile: Dockerfile
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@db:5432/${DB_NAME}?sslmode=disable
      PORT: 8080
      LOG_FORMAT: json
      APP_VERSION: ${APP_VERSION:-latest}
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped
    ports:
      - "8080:8080"
    command: ["./adhd-api"]
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  bot:
    build: 
      context: .
      dockerfile: Dockerfile
    environment:
      TELEGRAM_BOT_TOKEN: ${TELEGRAM_BOT_TOKEN}
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@db:5432/${DB_NAME}?sslmode=disable
      WEB_URL: ${WEB_URL}
      LOG_FORMAT: json
      APP_VERSION: ${APP_VERSION:-latest}
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped
    command: ["./adhd-bot"]

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        VITE_API_URL: ${API_URL}
        VITE_BOT_USERNAME: ${BOT_USERNAME}
    ports:
      - "3000:80"
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - api
      - frontend
    restart: unless-stopped

volumes:
  postgres_data:
```

**File**: `nginx.conf` (new)

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api {
        server api:8080;
    }
    
    upstream frontend {
        server frontend:80;
    }
    
    server {
        listen 80;
        server_name yourdomain.com;
        
        # Redirect HTTP to HTTPS
        return 301 https://$server_name$request_uri;
    }
    
    server {
        listen 443 ssl http2;
        server_name yourdomain.com;
        
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        
        # API routes
        location /api/ {
            proxy_pass http://api;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
        
        # Auth routes
        location /auth/ {
            proxy_pass http://api;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
        
        # Health check
        location /health {
            proxy_pass http://api;
        }
        
        # Frontend
        location / {
            proxy_pass http://frontend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

**File**: `.env.prod.example` (new)

```env
# Database
DB_NAME=adhd_bot_prod
DB_USER=postgres
DB_PASSWORD=your_secure_password

# Telegram
TELEGRAM_BOT_TOKEN=your_bot_token
BOT_USERNAME=YourBotUsername

# Web
WEB_URL=https://yourdomain.com
API_URL=https://yourdomain.com

# App
APP_VERSION=1.0.0
LOG_FORMAT=json
```

**Testing**:
- Test production deployment with Docker Compose
- Test SSL certificate configuration
- Test nginx routing
- Test health checks and monitoring

### 5.7 Documentation & User Guide
**File**: `docs/DEPLOYMENT.md` (new)

```markdown
# Deployment Guide

## Prerequisites

- Docker and Docker Compose
- Domain name with SSL certificate
- Telegram bot token
- PostgreSQL database (or use included container)

## Quick Start

1. Clone the repository
2. Copy environment file: `cp .env.prod.example .env.prod`
3. Update environment variables in `.env.prod`
4. Run: `docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d`

## Environment Variables