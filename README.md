<!-- @format -->

# Social Go API

[![Release Pipeline](https://github.com/yusuf-cirak/social-go/actions/workflows/release-please.yaml/badge.svg)](https://github.com/yusuf-cirak/social-go/actions/workflows/release-please.yaml)
[![CodeQL](https://github.com/yusuf-cirak/social-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/yusuf-cirak/social-go/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yusuf-cirak/social-go)](https://goreportcard.com/report/github.com/yusuf-cirak/social-go)
[![codecov](https://codecov.io/gh/yusuf-cirak/social-go/branch/master/graph/badge.svg)](https://codecov.io/gh/yusuf-cirak/social-go)
[![Docker Pulls](https://img.shields.io/docker/pulls/ghcr.io/yusuf-cirak/social-go)](https://ghcr.io/yusuf-cirak/social-go)
[![License](https://img.shields.io/github/license/yusuf-cirak/social-go)](LICENSE)

A modern social media API built with Go, featuring posts, comments, followers, and real-time interactions.

## Features

- 🚀 **High Performance**: Built with Go for optimal performance
- 🔐 **JWT Authentication**: Secure user authentication and authorization
- 📝 **Posts & Comments**: Full CRUD operations for social content
- 👥 **Follow System**: User following and follower relationships
- 🔄 **Real-time Feed**: Dynamic content feeds based on user interests
- 📊 **Pagination**: Efficient data pagination for large datasets
- 🏷️ **Tags**: Content tagging and filtering system
- 🐳 **Docker Ready**: Fully containerized with Docker support
- 📈 **Monitoring**: Health checks and observability features

## Tech Stack

- **Language**: Go 1.24.1
- **Framework**: Chi Router v5
- **Database**: PostgreSQL with migrations
- **Authentication**: JWT tokens
- **Logging**: Uber Zap
- **Validation**: Go Playground Validator
- **Testing**: Testify
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.24.1 or higher
- PostgreSQL 14+
- Docker (optional)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/yusuf-cirak/social-go.git
   cd social-go
   ```

2. **Set up environment variables**

   ```bash
   cp .envrc.example .envrc
   # Edit .envrc with your configuration
   source .envrc
   ```

3. **Install dependencies**

   ```bash
   go mod download
   ```

4. **Run database migrations**

   ```bash
   make migrate-up
   ```

5. **Start the application**
   ```bash
   go run cmd/api/main.go
   ```

### Using Docker

1. **Build and run with Docker Compose**

   ```bash
   docker-compose up --build
   ```

2. **Access the API**
   ```
   http://localhost:8080
   ```

## API Documentation

### Health Check

```bash
GET /v1/health
```

### Authentication

```bash
POST /v1/auth/register   # Register new user
POST /v1/auth/login      # Login user
POST /v1/auth/logout     # Logout user
```

### Users

```bash
GET    /v1/users/{id}    # Get user profile
PUT    /v1/users/{id}    # Update user profile
DELETE /v1/users/{id}    # Delete user account
```

### Posts

```bash
GET    /v1/posts         # Get posts feed
POST   /v1/posts         # Create new post
GET    /v1/posts/{id}    # Get specific post
PUT    /v1/posts/{id}    # Update post
DELETE /v1/posts/{id}    # Delete post
```

### Comments

```bash
GET    /v1/posts/{id}/comments  # Get post comments
POST   /v1/posts/{id}/comments  # Add comment
PUT    /v1/comments/{id}        # Update comment
DELETE /v1/comments/{id}        # Delete comment
```

### Followers

```bash
POST   /v1/users/{id}/follow    # Follow user
DELETE /v1/users/{id}/follow    # Unfollow user
GET    /v1/users/{id}/followers # Get followers
GET    /v1/users/{id}/following # Get following
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Database Operations

```bash
# Create new migration
make migrate-create name=your_migration_name

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down 1
```

### Code Quality

```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...

# Security scan
gosec ./...
```

## Database Schema

The application uses PostgreSQL with the following main tables:

- **users**: User accounts and profiles
- **posts**: User posts with content and metadata
- **comments**: Comments on posts
- **followers**: User follow relationships
- **tags**: Content tags and associations

See `cmd/migrate/migrations/` for detailed schema definitions.

## Configuration

Configuration is managed through environment variables:

```env
# Server
PORT=8080
HOST=localhost

# Database
DB_ADDR=postgres://user:password@localhost/dbname?sslmode=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

## Deployment

### Using Docker

```bash
# Build image
docker build -t social-go .

# Run container
docker run -p 8080:8080 social-go
```

### Using Docker Compose

```bash
docker-compose up -d
```

### Production Considerations

- Use PostgreSQL in production
- Set up proper SSL/TLS termination
- Configure rate limiting appropriately
- Set up monitoring and logging
- Use container orchestration (Kubernetes/Docker Swarm)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please read our [Contributing Guidelines](.github/pull_request_template.md) for details.

## Testing Strategy

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test API endpoints and database interactions
- **Security Tests**: Automated security scanning with CodeQL and Gosec
- **Performance Tests**: Load testing and benchmarks

## Monitoring & Observability

- Health check endpoints
- Structured logging with Zap
- Error tracking and reporting
- Performance metrics

## Security

- JWT-based authentication
- Rate limiting to prevent abuse
- Input validation and sanitization
- SQL injection prevention
- Security headers

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or need help, please:

1. Check the [Issues](https://github.com/yusuf-cirak/social-go/issues) page
2. Create a new issue using the appropriate template
3. Contact the maintainers

## Changelog

See [Releases](https://github.com/yusuf-cirak/social-go/releases) for a detailed changelog.

---

Made with ❤️ by [Yusuf Cirak](https://github.com/yusuf-cirak)
