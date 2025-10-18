# SpaceAPI Server - Go Project

This is a standalone Go-based server providing an endpoint for the SpaceAPI.

## Project Structure

```
spaceapi-endpoint/
├── cmd/
│   └── spaceapi/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/
│   │   └── spaceapi.go      # HTTP handlers for API endpoints
│   ├── middleware/
│   │   └── cors.go          # CORS middleware
│   ├── models/
│   │   └── spaceapi.go      # Data models and structures
│   └── services/
│       └── spaceapi.go      # Business logic and data loading
├── scripts/
│   ├── update-space-status.sh
│   └── update-people-count.sh
├── spaceapi.json            # Configuration file
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── Dockerfile.spaceapi      # Docker build configuration
├── docker-compose.yml       # Docker Compose configuration
├── Makefile                 # Build and development commands
└── README.md                # Main documentation
```

## Quick Start

### Development

```bash
# Install dependencies
make deps

# Build the application
make build

# Run the application
make run

# Or run directly
go run ./cmd/spaceapi
```

### Docker

```bash
# Build Docker image
make docker-build

# Run with Docker
make docker-run

# Or use docker-compose
make docker-compose-up
```

## API Endpoints

- `GET /api/space` - Get complete SpaceAPI JSON
- `POST /api/space/state` - Update space state (open/closed)
- `POST /api/space/people` - Update people count
- `POST /api/space/event` - Add an event
- `GET /health` - Health check

## Configuration

The server loads configuration from `spaceapi.json`. If the file is not found or invalid, it falls back to default values.

## Development Commands

```bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make clean          # Clean build artifacts
make fmt            # Format code
make lint           # Lint code
make docker-build   # Build Docker image
make docker-run     # Run Docker container
```

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router and URL matcher
- Go 1.21+

## License

see LICENSE.md
