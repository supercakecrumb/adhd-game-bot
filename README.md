# ADHD Game Bot 🎮

> Transform ADHD task management into an engaging game experience

[![Go Version](https://img.shields.io/badge/Go-1.25.0-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![Telegram](https://img.shields.io/badge/Telegram-Bot-blue.svg)](https://core.telegram.org/bots)

ADHD Game Bot is a sophisticated Telegram-based platform that gamifies daily task management for individuals with ADHD. By leveraging behavioral psychology and modern cloud architecture, we're making ADHD management engaging, rewarding, and sustainable.

## 🌟 Key Features

### For Users
- **🎯 Smart Task Management**: Categorized tasks (daily, weekly, ad-hoc) with intelligent scheduling
- **🏆 Gamification Engine**: Earn points, maintain streaks, unlock achievements
- **🛍️ Virtual Shop**: Spend earned points on motivating rewards
- **⏱️ Flexible Timers**: Countdown and stopwatch modes with pause/resume
- **🌍 Multi-timezone Support**: Works seamlessly across time zones
- **👥 Group Support**: Use in private chats or support groups

### For Developers
- **🏗️ Clean Architecture**: Domain-driven design with clear separation of concerns
- **🔄 RESTful API**: Complete API for external integrations
- **🔒 Idempotent Operations**: Safe transaction handling
- **📊 Comprehensive Testing**: Unit and integration test coverage
- **🐳 Docker Ready**: Easy deployment with Docker Compose
- **📈 Scalable Design**: Built to handle millions of users

## 📚 Documentation

- **[Investor Pitch](INVESTOR_PITCH.md)** - Business overview and investment opportunity
- **[Technical Overview](TECHNICAL_OVERVIEW.md)** - Detailed technical documentation
- **[User Guide](USER_GUIDE.md)** - Complete guide for end users
- **[API Documentation](#api-endpoints)** - REST API reference

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- A Telegram bot token from [@BotFather](https://t.me/BotFather)
- PostgreSQL 13+ (if running locally)
- Go 1.25.0+ (for local development)

### Running with Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd adhd-game-bot
   ```

2. Set your Telegram bot token as an environment variable:
   ```bash
   export TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
   ```

   For testing purposes, you can use a dummy token (note that the bot won't be able to connect to Telegram):
   ```bash
   export TELEGRAM_BOT_TOKEN=dummy_token_for_testing
   ```

3. Build and start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

   If you encounter warnings about orphan containers, you can clean them up:
   ```bash
   docker-compose up --build --remove-orphans
   ```

4. The bot should now be running and connected to Telegram (if using a valid token).
   The API service will be available at http://localhost:8080

## 🗄️ Database Setup

### Running Migrations

Before first use, apply database migrations:

```bash
# Using Docker Compose
docker-compose run --rm bot ./migrate

# Or locally
./scripts/migrate.sh
```

Migrations are located in `internal/infra/postgres/migrations/` and run sequentially.

## 💻 Development

### Local Development Setup

1. **Install Go 1.25.0+**
   ```bash
   # Check version
   go version
   ```

2. **Clone and Install Dependencies**
   ```bash
   git clone <repository-url>
   cd adhd-game-bot
   go mod download
   ```

3. **Set Environment Variables**
   ```bash
   export TELEGRAM_BOT_TOKEN=your_token_here
   export DATABASE_URL=postgres://user:pass@localhost/adhd_bot
   ```

4. **Build and Run**
   ```bash
   # Build
   go build -o adhd-bot cmd/bot/main.go
   go build -o adhd-api cmd/api/main.go
   
   # Run
   ./adhd-bot  # Terminal 1
   ./adhd-api  # Terminal 2
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/usecase/...
```

## 🤖 Bot Commands

### Available Commands
- `/start` - Initialize the bot and create your profile
- `/shop` - Browse available rewards
- `/buy <item_code>` - Purchase items with earned points
- `/balance` - Check your current point balance
- `/help` - Get command list and assistance

### Coming Soon
- `/tasks` - View and manage your tasks
- `/complete <task_id>` - Mark tasks as complete
- `/streak` - View your habit streaks
- `/settings` - Configure preferences

## 📡 API Endpoints

The REST API enables external integrations and custom clients.

### Base URL
```
http://localhost:8080/api
```

### Task Management

#### Create Task
```http
POST /api/tasks?user_id={user_id}
Content-Type: application/json

{
  "chat_id": 123456789,
  "title": "Morning Meditation",
  "description": "10 minutes of mindfulness",
  "category": "daily",
  "difficulty": "easy",
  "schedule_json": "{}",
  "base_duration": 600,
  "grace_period": 300,
  "cooldown": 0,
  "reward_curve_json": "{}",
  "streak_enabled": true,
  "status": "active",
  "time_zone": "America/New_York"
}
```

#### Get Task
```http
GET /api/tasks/{task_id}
```

#### Update Task
```http
PUT /api/tasks/{task_id}
Content-Type: application/json

{
  "title": "Updated Task Title",
  "difficulty": "medium"
}
```

#### Complete Task
```http
POST /api/tasks/{task_id}/complete?user_id={user_id}
```

#### List User Tasks
```http
GET /api/users/{user_id}/tasks
```

### Response Format
```json
{
  "data": { ... },
  "error": null,
  "timestamp": "2025-01-09T15:00:00Z"
}
```

## 🏗️ Architecture

The project implements Clean Architecture principles:

```
├── cmd/                    # Application entry points
│   ├── bot/               # Telegram bot executable
│   └── api/               # REST API server
├── internal/              # Private application code
│   ├── domain/            # Core business logic
│   │   ├── entity/        # Domain entities
│   │   └── valueobject/   # Value objects
│   ├── usecase/           # Application use cases
│   ├── infra/             # Infrastructure implementations
│   │   ├── postgres/      # PostgreSQL repositories
│   │   └── http/          # HTTP handlers
│   └── ports/             # Interface definitions
├── test/                  # Test files
│   ├── acceptance/        # End-to-end tests
│   └── fixtures/          # Test data builders
└── scripts/               # Utility scripts
```

### Key Design Patterns
- **Repository Pattern**: Abstract data access
- **Dependency Injection**: Loose coupling
- **Domain-Driven Design**: Rich domain models
- **CQRS**: Separate read/write operations
- **Event Sourcing**: Audit trail (planned)

## 🔒 Security

- **Data Encryption**: All sensitive data encrypted at rest
- **TLS Communication**: Secure API endpoints
- **Input Validation**: Comprehensive sanitization
- **SQL Injection Prevention**: Parameterized queries
- **Rate Limiting**: API abuse prevention
- **Audit Logging**: Security event tracking

## 🚀 Deployment

### Production Deployment

1. **Environment Setup**
   ```bash
   cp .env.example .env
   # Edit .env with production values
   ```

2. **Deploy with Docker Compose**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Configure Monitoring**
   - Set up Prometheus metrics collection
   - Configure Grafana dashboards
   - Enable alerting

### Scaling Considerations
- Horizontal scaling via Kubernetes
- Database read replicas for performance
- Redis caching layer
- CDN for static assets

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Process
1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Code Standards
- Follow Go best practices
- Maintain test coverage above 80%
- Document all public APIs
- Use meaningful commit messages

## 📊 Project Status

- ✅ Core bot functionality
- ✅ Shop system
- ✅ REST API
- ✅ Database migrations
- ✅ Docker support
- 🚧 Task notifications
- 🚧 Mobile apps
- 📅 Analytics dashboard
- 📅 Third-party integrations

## 🆘 Support

- **Documentation**: See our [comprehensive guides](USER_GUIDE.md)
- **Issues**: Report bugs via [GitHub Issues](https://github.com/your-repo/issues)
- **Community**: Join our [Telegram group](https://t.me/ADHDGameBotCommunity)
- **Email**: support@adhd-game-bot.com

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- The ADHD community for invaluable feedback
- Contributors who've helped shape this project
- Open source libraries that made this possible

---

**Built with ❤️ for the ADHD community**
