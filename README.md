# ADHD Game Bot

A Telegram bot designed to help individuals with ADHD manage their daily routines through gamification.

## Features

- Task management with rewards
- Shop system for purchasing items
- Balance tracking
- Multi-chat support

## Prerequisites

- Docker and Docker Compose
- A Telegram bot token (obtained from [@BotFather](https://t.me/BotFather))

## Building and Running

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

## Database Migrations

Before the bot can function properly, you need to apply the database migrations:

```bash
# Run migrations using docker-compose
docker-compose run --rm bot ./migrate
```

The `--rm` flag ensures that the container is removed after it finishes running, preventing orphan containers.

The migration tool will apply all pending migrations to set up the database schema. After running the migrations, you can start the bot normally:
```bash
docker-compose up -d
```

The migration files are located in `internal/infra/postgres/migrations/` and follow a sequential numbering scheme.

## Development

To build the bot locally without Docker:

1. Ensure you have Go 1.25.0+ installed
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Build the bot:
   ```bash
   go build -o adhd-bot cmd/bot/main.go
   ```
4. Build the API:
   ```bash
   go build -o adhd-api cmd/api/main.go
   ```
5. Run the bot:
   ```bash
   ./adhd-bot
   ```
6. Run the API (in a separate terminal):
   ```bash
   ./adhd-api
   ```

Note: You'll need to set up a PostgreSQL database separately for local development.

## Bot Commands

- `/start` - Start the bot and register as a new user
- `/shop` - View available items in the shop
- `/buy <item_code>` - Purchase an item from the shop
- `/balance` - Check your current balance

## API Endpoints

The API service provides RESTful endpoints for task management:

### Tasks

- `POST /api/tasks?user_id={user_id}` - Create a new task
- `GET /api/tasks/{task_id}` - Get a specific task
- `PUT /api/tasks/{task_id}` - Update a task
- `POST /api/tasks/{task_id}/complete?user_id={user_id}` - Mark a task as complete

### User Tasks

- `GET /api/users/{user_id}/tasks` - List all tasks for a user

### Task Data Structure

When creating or updating tasks, use the following JSON structure:

```json
{
  "chat_id": 123456789,
  "title": "Example Task",
  "description": "This is an example task",
  "category": "daily",
  "difficulty": "medium",
  "schedule_json": "{}",
  "base_duration": 3600,
  "grace_period": 300,
  "cooldown": 0,
  "reward_curve_json": "{}",
  "streak_enabled": true,
  "status": "active",
  "time_zone": "UTC"
}
```

## Architecture

The bot follows a clean architecture pattern with:
- Domain layer containing entities and value objects
- Use case layer with business logic
- Infrastructure layer with database implementations
- Ports and adapters for dependency inversion

## License

This project is licensed under the MIT License.
