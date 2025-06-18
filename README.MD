# Nucleus - Event-Driven E-Commerce Backend

A production-grade backend system built with Go, following Clean Architecture principles and enterprise-level best practices.

## 🏗️ Architecture

- **Clean Architecture**: Separation of concerns with clearly defined layers
- **Event-Driven**: Using RabbitMQ for commands and Kafka for events
- **Microservices Ready**: Designed for scalability and maintainability
- **Cloud Native**: Containerized with Docker and Kubernetes ready

## 🛠️ Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **Event Streaming**: Kafka
- **Authentication**: Keycloak (OAuth2/OIDC)
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logs with Loki
- **API Documentation**: Swagger/OpenAPI
- **Container**: Docker & Docker Compose

## 📁 Project Structure

```
nucleus/
├── cmd/                    # Application entrypoints
│   ├── api/               # REST API server
│   └── worker/            # Background workers
├── internal/              # Private application code
│   ├── domain/            # Enterprise business rules
│   ├── application/       # Application business rules
│   ├── infrastructure/    # Frameworks and drivers
│   └── interfaces/        # Interface adapters
├── pkg/                   # Public packages
├── configs/               # Configuration files
├── deployments/           # Deployment configurations
├── scripts/               # Utility scripts
├── tests/                 # Test suites
└── docs/                  # Documentation
```

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- Make (optional but recommended)

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/vashafaris/nucleus.git
cd nucleus
```

2. Set up the environment:
```bash
make setup
```

3. Start all services:
```bash
make docker-up
```

4. Run the application:
```bash
make run
```

The API will be available at http://localhost:8080

### Available Services

- **API**: http://localhost:8080
- **Keycloak**: http://localhost:8180
- **RabbitMQ Management**: http://localhost:15672 (user: nucleus, password: rabbitmq123)
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (user: admin, password: admin123)

## 📝 Development

### Common Commands

```bash
# Install dependencies
make deps

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Generate mocks
make mocks

# Build the application
make build

# Generate Swagger docs
make swagger
```

### Running Locally

```bash
# Copy environment variables
cp .env.example .env

# Start dependencies
make docker-up

# Run the application
make run
```

### Database Migrations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Seed database
make seed
```

## 🧪 Testing

The project includes three levels of testing:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **E2E Tests**: Test complete user flows

Run all tests:
```bash
make test
```

Run with coverage:
```bash
make test-coverage
```

## 📊 Monitoring

### Metrics
- Prometheus metrics available at `/metrics`
- Pre-configured Grafana dashboards
- Application metrics: request latency, error rates, etc.

### Logging
- Structured JSON logging
- Centralized log aggregation with Loki
- Log levels: debug, info, warn, error

### Health Checks
- Liveness: `/health/live`
- Readiness: `/health/ready`

## 🔒 Security

- OAuth2/OIDC authentication with Keycloak
- JWT token validation
- Role-based access control (RBAC)
- Rate limiting
- CORS configuration
- Security headers

## 📚 API Documentation

Swagger documentation is available at:
- http://localhost:8080/swagger/index.html

To regenerate docs:
```bash
make swagger
```

## 🚢 Deployment

### Docker

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-up
```

### Kubernetes

Kubernetes manifests are available in `deployments/k8s/`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👥 Author

- GitHub: [@vashafaris](https://github.com/vashafaris)