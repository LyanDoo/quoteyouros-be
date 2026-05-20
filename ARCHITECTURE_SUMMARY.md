# QuoteYourOS Backend - Architecture & Tech Summary

## 🎯 Executive Summary

A production-ready REST API backend for QuoteYourOS built with **Go**, **PostgreSQL**, and **Fiber** following **Clean Architecture** principles. The system is fully containerized with Docker and ready for development and deployment.

---

## 🏗️ Architecture Overview

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────┐
│                  HTTP/REST API                   │
│              (Fiber Framework)                   │
├─────────────────────────────────────────────────┤
│              Presentation Layer                  │
│    - HTTP Handlers (.Handler)                    │
│    - Request/Response DTOs                       │
│    - Route Configuration                         │
├─────────────────────────────────────────────────┤
│             Business Logic Layer                 │
│    - Use Cases (implements Domain interfaces)    │
│    - Validation & transformation                 │
│    - Error handling                              │
├─────────────────────────────────────────────────┤
│              Domain Layer                        │
│    - Entities (Blog, Project, User, etc.)        │
│    - Repository Interfaces                       │
│    - UseCase Interfaces                          │
├─────────────────────────────────────────────────┤
│           Infrastructure Layer                   │
│    - PostgreSQL repositories                     │
│    - Email service (Resend/SMTP/NoOp)            │
│    - File storage (Local/S3)                     │
│    - Database connection pool                    │
├─────────────────────────────────────────────────┤
│             External Services                    │
│    - PostgreSQL Database                         │
│    - Email APIs (Resend, SendGrid, SMTP)         │
│    - AWS S3 (optional file storage)              │
└─────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Dependency Inversion** - High-level modules depend on abstractions (interfaces), not concrete implementations
2. **Single Responsibility** - Each layer/package has one reason to change
3. **Open/Closed** - Easy to extend (e.g., adding new email providers) without modifying existing code
4. **Interface Segregation** - Repository interfaces are focused and specific
5. **Testability** - All dependencies are injected, making mocking easy

---

## 📊 Technology Stack

### Core Framework
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Language** | Go 1.21+ | Fast, compiled, excellent concurrency support |
| **Web Framework** | Fiber v2 | High performance, minimal overhead, similar to Express.js |
| **Database** | PostgreSQL 16 | Robust, ACID compliance, excellent JSON support |
| **Database Driver** | pgx v5 | Type-safe, high performance, connection pooling |

### Authentication & Security
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **JWT** | golang-jwt/jwt/v5 | Industry standard, stateless auth |
| **Password Hashing** | golang.org/x/crypto | Battle-tested, secure algorithms |
| **CORS** | Fiber middleware | Built-in support |

### Configuration & Logging
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Config Management** | spf13/viper | Flexible, supports multiple formats, env vars |
| **Structured Logging** | uber/zap | High performance, structured output |

### Infrastructure
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Containerization** | Docker | Consistent dev/prod environments |
| **Orchestration** | Docker Compose | Simple multi-service setup |
| **Email (Optional)** | Resend/SendGrid/SMTP | Flexible provider support |
| **File Storage (Optional)** | AWS S3/Local | Scalable or local development |

---

## 📁 Project Structure Explained

```
quoteyouros-be/
│
├── cmd/
│   └── main.go                      # Application entry point
│                                     # - Initializes services
│                                     # - Sets up routes
│                                     # - Starts server
│
├── internal/                         # Private application code (not importable)
│   │
│   ├── config/
│   │   └── config.go               # Configuration loader using Viper
│   │
│   ├── domain/                      # **DOMAIN LAYER** - Core business logic
│   │   ├── entities.go             # Domain models + DTOs
│   │   ├── repositories.go         # Repository interfaces
│   │   └── usecases.go             # UseCase interfaces
│   │
│   ├── usecase/                     # **BUSINESS LOGIC LAYER** (NOT YET IMPLEMENTED)
│   │   ├── blog/                   # Blog use cases
│   │   ├── project/                # Project use cases
│   │   ├── contact/                # Contact use cases
│   │   └── auth/                   # Authentication use cases
│   │
│   ├── handler/                     # **PRESENTATION LAYER** (NOT YET IMPLEMENTED)
│   │   ├── blog.go                 # Blog HTTP handlers
│   │   ├── project.go              # Project HTTP handlers
│   │   ├── contact.go              # Contact HTTP handlers
│   │   ├── auth.go                 # Authentication handlers
│   │   └── profile.go              # Profile handlers
│   │
│   ├── middleware/                  # HTTP middleware
│   │   ├── auth.go                 # JWT authentication middleware
│   │   └── cors.go                 # CORS configuration
│   │
│   ├── repository/                  # **INFRASTRUCTURE LAYER** - Data Access
│   │   ├── blog/postgres.go        # Blog repository implementation
│   │   ├── project/postgres.go     # Project repository implementation
│   │   ├── contact/postgres.go     # Contact repository implementation
│   │   ├── message/postgres.go     # Message repository implementation
│   │   └── user/postgres.go        # User repository implementation
│   │
│   └── infrastructure/
│       ├── postgres/
│       │   └── connection.go       # PostgreSQL connection pool
│       ├── email/
│       │   └── service.go          # Email service (Resend/SMTP/NoOp)
│       └── file/
│           └── service.go          # File storage service (Local/S3)
│
├── pkg/                             # **Shared Utilities** (reusable across packages)
│   ├── errors/
│   │   └── errors.go               # Custom error types & constructors
│   └── response/
│       └── response.go             # Standard response formatters
│
├── migrations/
│   └── 001_init_schema.sql         # Database schema with indexes
│
├── docker-compose.yml              # Multi-container setup (PostgreSQL + Backend)
├── Dockerfile                      # Container build configuration
├── .env.example                    # Environment variables template
├── go.mod                          # Go module dependencies
├── go.sum                          # Dependency checksums
├── Makefile                        # Development commands
├── README.md                       # User-facing documentation
├── DEVELOPMENT_GUIDE.md            # Developer implementation guide
└── ARCHITECTURE_SUMMARY.md         # This file
```

### Layer Responsibilities

**Domain Layer** (`internal/domain/`)
- Pure business logic, independent of frameworks
- Entities: Blog, Project, User, Contact, Message
- Interfaces: Repository contracts, UseCase contracts
- No external dependencies

**UseCase Layer** (`internal/usecase/*/`)
- Implements domain interfaces
- Orchestrates business logic
- Handles validation and transformation
- Depends on domain layer only

**Handler Layer** (`internal/handler/`)
- HTTP request/response handling
- Route setup
- Calls use cases
- Returns formatted responses

**Repository Layer** (`internal/repository/*/`)
- Implements domain repository interfaces
- Database queries
- Depends only on domain entities

**Middleware Layer** (`internal/middleware/`)
- Cross-cutting concerns
- JWT validation
- CORS handling
- Error handling

---

## 🔄 Request Flow Example

### Public Blog List Request
```
GET /api/blog?page=1&limit=10

1. Request arrives at Fiber router
2. Route handler: handler.BlogHandler.GetAllBlogPosts()
3. Extract pagination params from query
4. Call usecase: BlogUseCase.GetAllBlogPosts(ctx, page, limit)
5. UseCase calls repository: BlogRepository.GetAllBlogPosts(ctx, limit, offset)
6. Repository executes SQL query against PostgreSQL
7. Returns []*BlogPost, total count
8. UseCase processes results
9. Handler formats PaginatedResponse
10. Return JSON response to client
```

### Protected Blog Create Request
```
POST /api/blog

1. Request arrives with Authorization header
2. JWTAuth middleware validates token
3. Extract user info from token (stored in c.Locals)
4. Route handler: handler.BlogHandler.CreateBlogPost()
5. Parse and validate request body
6. Call usecase: BlogUseCase.CreateBlogPost(ctx, req)
7. UseCase creates domain entity
8. UseCase calls repository: BlogRepository.CreateBlogPost(ctx, post)
9. Repository inserts into database
10. Handler returns created post with 201 status
```

---

## 🔐 Security Measures

1. **JWT Authentication** - All protected endpoints require valid JWT token
2. **Password Hashing** - Passwords hashed with bcrypt (via crypto package)
3. **Environment Variables** - Secrets not hardcoded, loaded from .env
4. **CORS Configuration** - Whitelist allowed origins
5. **SQL Injection Prevention** - All queries use parameterized statements
6. **Error Handling** - Generic error messages to avoid information leakage

---

## 📈 Scalability Considerations

### Current Implementation
- ✅ Connection pooling (pgxpool)
- ✅ Pagination for list endpoints
- ✅ Database indexes on common queries
- ✅ Containerized for easy deployment

### Future Improvements
- [ ] Caching layer (Redis)
- [ ] Rate limiting
- [ ] Request/response compression
- [ ] API versioning strategy
- [ ] GraphQL support (optional)
- [ ] Message queue (Kafka/RabbitMQ) for async tasks
- [ ] Database sharding/replication
- [ ] CDN for static files

---

## 🧪 Testing Strategy

### Unit Tests
- Test use cases with mock repositories
- Test handlers with mock use cases
- Location: Alongside implementation files (`*_test.go`)

### Integration Tests
- Test full request flow with real database
- Use Docker Compose for test database
- Location: `tests/integration/` (to be created)

### API Tests
- Test endpoints with curl/Postman
- Use Postman collections or Thunder Client

---

## 🚀 Deployment

### Development
```bash
make db-up      # Start PostgreSQL
go run ./cmd    # Start API server
```

### Docker
```bash
docker-compose up       # Start all services
docker-compose down     # Stop all services
```

### Production (on server/cloud)
```bash
# Build image
docker build -t quoteyouros-backend:latest .

# Run with proper .env configuration
docker run -d \
  -p 8000:8000 \
  --env-file .env.prod \
  quoteyouros-backend:latest
```

---

## 📚 API Documentation

### Response Format

**Success Response**
```json
{
  "success": true,
  "data": { /* entity data */ },
  "message": "Operation completed successfully"
}
```

**Paginated Response**
```json
{
  "success": true,
  "data": [ /* array of entities */ ],
  "page": 1,
  "limit": 10,
  "total": 50,
  "message": "Data retrieved successfully"
}
```

**Error Response**
```json
{
  "success": false,
  "error": "Error message",
  "code": 400
}
```

### Status Codes
- `200 OK` - Successful GET/PUT
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing/invalid token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Server error

---

## 🎯 Next Implementation Steps

1. **Implement Use Cases** (blog, project, contact, auth)
2. **Create HTTP Handlers** for all endpoints
3. **Add Input Validation** middleware
4. **Complete Email Service** (Resend/SMTP)
5. **Implement File Upload** for resume
6. **Write Tests** (unit + integration)
7. **Add Logging** throughout application
8. **Performance Testing** and optimization
9. **Documentation** (API docs, OpenAPI spec)

---

## 📞 Quick Reference

| Task | Command |
|------|---------|
| Start dev | `make run` |
| Start DB | `make db-up` |
| Run tests | `make test` |
| Build binary | `make build` |
| Format code | `make fmt` |
| View logs | `make docker-logs` |
| Migrate DB | `make db-migrate` |

---

## 📖 Resources

- [Fiber Docs](https://docs.gofiber.io)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8949)
- [Go Best Practices](https://golang.org/doc/effective_go)

---

**Built with Clean Architecture principles for maintainability, testability, and scalability. 🚀**
