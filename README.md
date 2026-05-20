# QuoteYourOS Backend API

A clean architecture REST API backend built with Go, PostgreSQL, and Fiber framework.

## 📋 Project Structure

```
quoteyouros-be/
├── cmd/                          # Application entry point
│   └── main.go
├── internal/                      # Private application code
│   ├── config/                   # Configuration management
│   ├── domain/                   # Domain entities and interfaces
│   ├── usecase/                  # Business logic
│   ├── handler/                  # HTTP handlers
│   ├── middleware/               # HTTP middleware
│   ├── repository/               # Data access layer
│   └── infrastructure/           # External services (DB, Email, etc.)
├── pkg/                          # Shared packages
│   ├── errors/                   # Error handling
│   └── response/                 # Response formatting
├── migrations/                   # Database migrations
├── docker-compose.yml            # Docker services
├── Dockerfile                    # Container build config
├── go.mod                        # Go module definition
└── .env.example                  # Environment variables template
```

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

- **Domain Layer**: Business entities and interfaces (repository & usecase)
- **UseCase Layer**: Business logic implementation
- **Handler Layer**: HTTP request/response handling
- **Repository Layer**: Database operations
- **Infrastructure Layer**: External integrations (PostgreSQL, Email, etc.)

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 16+ (or use Docker)

### Setup

1. **Clone and navigate to project**
   ```bash
   cd quoteyouros-be
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Start PostgreSQL with Docker**
   ```bash
   docker-compose up -d postgres
   ```

4. **Run database migrations**
   ```bash
   # Using migrate tool or manual SQL execution
   psql -U postgres -d quoteyouros -f migrations/001_init_schema.sql
   ```

5. **Download Go dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

6. **Run the application**
   ```bash
   go run ./cmd
   ```

The API will be available at `http://localhost:8000`

## 📡 API Endpoints

### Public Endpoints

#### Blog
- `GET /api/blog` - List all blog posts (paginated)
- `GET /api/blog/:id` - Get a specific blog post

#### Projects
- `GET /api/projects` - List all projects
- `GET /api/projects/:id` - Get a specific project

#### Contact
- `POST /api/contact` - Submit contact form

#### Profile
- `GET /api/profile/about` - Get about me text
- `GET /api/profile/resume` - Get resume data
- `GET /api/profile/resume/download` - Download resume PDF

#### Authentication
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/auth/me` - Get current user info (requires auth)

### Protected Endpoints (Admin)

#### Blog Management
- `POST /api/blog` - Create blog post
- `PUT /api/blog/:id` - Update blog post
- `DELETE /api/blog/:id` - Delete blog post

#### Project Management
- `POST /api/projects` - Create project
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project

#### Messages
- `GET /api/messages` - Get all contact messages
- `DELETE /api/messages/:id` - Delete a message

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Protected endpoints require an `Authorization` header:

```
Authorization: Bearer <token>
```

## 📦 Dependencies

- **gofiber/fiber/v2** - Web framework
- **jackc/pgx/v5** - PostgreSQL driver
- **golang-jwt/jwt/v5** - JWT handling
- **spf13/viper** - Configuration management
- **go.uber.org/zap** - Structured logging
- **golang.org/x/crypto** - Password hashing

## 🐳 Docker

### Development

```bash
# Start all services
docker-compose up

# Stop all services
docker-compose down

# View logs
docker-compose logs -f backend
```

### Production

Update the `Dockerfile` and `.env` for production settings before building:

```bash
docker build -t quoteyouros-backend:latest .
docker run -p 8000:8000 --env-file .env quoteyouros-backend:latest
```

## 🔧 Development

### Make Commands

```bash
make build      # Build the binary
make run        # Run the application
make test       # Run tests
make db-up      # Start database
make db-down    # Stop database
make db-reset   # Reset database
```

### Adding New Features

1. Define domain entities in `internal/domain/entities.go`
2. Create repository interface in `internal/domain/repositories.go`
3. Implement repository in `internal/repository/<feature>/`
4. Create usecase interface in `internal/domain/usecases.go`
5. Implement usecase in `internal/usecase/<feature>/`
6. Create HTTP handler in `internal/handler/<feature>.go`
7. Register routes in `cmd/main.go`

## 📝 Environment Variables

See `.env.example` for all available configuration options.

Key variables:
- `SERVER_PORT` - API port (default: 8000)
- `DB_*` - Database connection details
- `JWT_SECRET` - JWT signing secret
- `JWT_EXPIRATION` - Token expiration time
- `EMAIL_PROVIDER` - Email service (resend, sendgrid, smtp)
- `CORS_ORIGINS` - Allowed CORS origins

## 🧪 Testing

```bash
go test ./...
go test -v ./...
go test -cover ./...
```

## 📚 Additional Resources

- [Fiber Documentation](https://docs.gofiber.io)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT Introduction](https://jwt.io)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📄 License

MIT

## 👨‍💻 Author

Built with ❤️ for QuoteYourOS
