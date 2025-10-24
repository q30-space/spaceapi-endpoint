# SpaceAPI Endpoint

[![CI](https://github.com/q30-space/spaceapi-endpoint/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/q30-space/spaceapi-endpoint/actions/workflows/ci.yml)

An API server for a [SpaceAPI](https://spaceapi.io/) endpoint.

## Overview

SpaceAPI-Endpoint provides a RESTful API that follows the [SpaceAPI v15 specification](https://spaceapi.io/docs/) to expose real-time information about the hackerspace status, plus a CLI tool for status monitoring.

## Architecture

- **SpaceAPI Server**: Go-based REST API server running on port 8089
- **SpaceIcon CLI**: Command-line tool for displaying space status with icons
- **Configuration**: JSON-based configuration file (`spaceapi.json`)

## Services

### 1. SpaceAPI Server (spaceapi)
- **Port**: 8089 (configurable via PORT env var, defaults to 8080)
- **Purpose**: Provides SpaceAPI endpoints
- **Health Check**: `/health`
- **Binary**: `bin/spaceapi`

### 2. SpaceIcon CLI (spaceicon)
- **Purpose**: Command-line tool for monitoring the space status of your endpoint from your terminal
- **Features**: Colored terminal output and i3blocks integration
- **Binary**: `bin/spaceicon`

## API Endpoints

### GET `/api/space`
Returns the complete SpaceAPI JSON response.

**Example:**
```bash
curl http://localhost:8089/api/space
```

### POST `/api/space/state` 🔒
Updates the space state (open/closed status). **Requires API key authentication.**

**Headers:**
```
X-API-Key: your_api_key_here
Content-Type: application/json
```

**Payload:**
```json
{
    "open": true,
    "message": "Space is open for members",
    "trigger_person": "John Doe"
}
```

**Example:**
```bash
curl -X POST \
  -H "X-API-Key: your_api_key_here" \
  -H "Content-Type: application/json" \
  -d '{"open": true, "message": "Space is open"}' \
  http://localhost:8089/api/space/state
```

### POST `/api/space/people` 🔒
Updates the people count in the space. **Requires API key authentication.**

**Headers:**
```
X-API-Key: your_api_key_here
Content-Type: application/json
```

**Payload:**
```json
{
    "value": 5,
    "location": "Main Space"
}
```

**Example:**
```bash
curl -X POST \
  -H "X-API-Key: your_api_key_here" \
  -H "Content-Type: application/json" \
  -d '{"value": 5, "location": "Main Space"}' \
  http://localhost:8089/api/space/people
```

### POST `/api/space/event` 🔒
Adds an event to the space timeline. **Requires API key authentication.**

**Headers:**
```
X-API-Key: your_api_key_here
Content-Type: application/json
```

**Payload:**
```json
{
    "name": "John Doe",
    "type": "check-in",
    "extra": "Working on Arduino project"
}
```

**Example:**
```bash
curl -X POST \
  -H "X-API-Key: your_api_key_here" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "type": "check-in"}' \
  http://localhost:8089/api/space/event
```

## Authentication & Rate Limiting

### API Key Authentication
All POST endpoints require authentication via API key. 

Check the Configuration section below.

### Rate Limiting
Failed authentication attempts are rate limited to mitigate brute force attacks:

- **Limit**: 5 failed attempts within 15 minutes
- **Block Duration**: 1 hour
- **Response**: HTTP 429 Too Many Requests with `Retry-After` header
- **Scope**: Per IP address

### Error Responses

| Status | Description | Response |
|--------|-------------|----------|
| 401 | Missing or invalid API key | `"API key required"` or `"Invalid API key"` |
| 403 | API key format invalid | `"Access forbidden"` |
| 429 | Rate limited | `"Too many failed authentication attempts. Please try again later."` |

## Configuration

### SpaceAPI Configuration
Copy `spaceapi.json.example` to `spaceapi.json`.

The SpaceAPI configuration is stored in `spaceapi.json`. Key sections:

- **Basic Info**: Space name, logo, URL
- **Location**: Address, coordinates, timezone
- **Contact**: Email, IRC, social media
- **State**: Current open/closed status
- **Sensors**: People count, temperature, etc.
- **Feeds**: Blog RSS, calendar feeds

For full documentation check the [Schema Documentation](https://spaceapi.io/docs/) .

### Authentication Setup

1. **Copy the environment template**:
   ```bash
   cp .env.example .env
   ```

2. **Generate a secure API key**:
   ```bash
   openssl rand -hex 32
   ```

3. **Set your API key in `.env`**:
   ```bash
   # Edit .env file
   SPACEAPI_AUTH_KEY=your_generated_key_here
   ```

4. **For Docker Compose**:
   The `.env` file is automatically loaded by Docker Compose.

5. **For manual scripts**:
   ```bash
   export SPACEAPI_AUTH_KEY=your_generated_key_here
   ./scripts/update-space-status.sh open
   ```

## CLI Tools

### SpaceIcon CLI

The `spaceicon` tool provides a simple way to check space status with visual indicators.

#### Basic Usage
```bash
# Check space status (colored output)
./bin/spaceicon https://your-spaceapi-url/api/space

# i3blocks integration (three-line output)
./bin/spaceicon --i3block https://your-spaceapi-url/api/space
```

#### Output Modes

**Default Mode** (colored terminal output):
- Space open: Green icon
- Space closed: Red icon
- Error: Red icon

**i3blocks Mode** (`--i3block` flag):
- Line 1: Icon
- Line 2: Extended title
- Line 3: Hex color
  - Open: `#228800`
  - Closed/Error: `#FF0F0F`

#### Examples
```bash
# Local development
./bin/spaceicon http://localhost:8089/api/space

# Production API
./bin/spaceicon https://localhost:8089/api/space

# i3blocks configuration
./bin/spaceicon --i3block https://localhost:8089/api/space
```

## Management Scripts

### Update Space Status
```bash
# Open the space
./scripts/update-space-status.sh open "Open for members" "John Doe"

# Close the space
./scripts/update-space-status.sh closed "Closed for maintenance" "Jane Smith"
```

### Update People Count
```bash
# Set people count
./scripts/update-people-count.sh 5 "Main Space"

# Set to zero
./scripts/update-people-count.sh 0
```

### Space Status Icon (Bash)
```bash
# Simple status check with icons
./scripts/space_status_icon.sh https://your-spaceapi-url/api/space
```

## Testing

The project includes a comprehensive test suite covering all critical components.

### Running Tests

```bash
# Run all tests (recommended)
make test

# Run with verbose output
make test-verbose

# Run with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

Or run directly with Go:
```bash
# Run all tests
go test ./internal/... ./cmd/...

# Run tests with verbose output
go test -v ./internal/... ./cmd/...

# Run tests with coverage report
go test -cover ./internal/... ./cmd/...

# Run tests for specific packages
go test ./internal/handlers ./internal/middleware

# Run tests with detailed coverage
go test -coverprofile=coverage.out ./internal/handlers ./internal/middleware
go tool cover -html=coverage.out
```

### Test Coverage

The test suite provides:
- **96.1% code coverage** for handlers
- **Integration tests** for all HTTP endpoints
- **CORS middleware** tests (10 tests)
- **Mocked data** (no external file dependencies)

### Test Structure

```
internal/
├── handlers/
│   └── spaceapi_test.go      # HTTP endpoint tests
├── middleware/
│   └── cors_test.go          # CORS middleware tests
└── testutil/
    └── helpers.go            # Test data helpers
```

### Test Features

- **HTTP Integration Tests**: Uses `httptest` for realistic HTTP testing
- **Mocked Data**: No need for `spaceapi.json` file during testing
- **Testify Suite**: Organized test structure with setup/teardown
- **Concurrent Testing**: Middleware concurrency tests
- **Comprehensive Coverage**: All critical paths and edge cases

## Building

### Build All Binaries
```bash
# Build both spaceapi server and spaceicon CLI
make build

# Or build individually
go build -o bin/spaceapi ./cmd/spaceapi
go build -o bin/spaceicon ./cmd/spaceicon
```

### Development
```bash
# Run the SpaceAPI server
make run
# or
go run ./cmd/spaceapi

# Test the CLI tool
./bin/spaceicon http://localhost:8089/api/space
```

## Deployment

### Production - Docker Image (Recommended) 🐳

The easiest way to deploy is using the pre-built Docker image from GitHub Container Registry.

📖 **For a complete guide, see the [Deployment Guide](doc/DEPLOYMENT_GUIDE.md).**

> **⚠️ First-Time Setup**: The Docker image needs to be built by CI first. If you get "unauthorized" or "not found" errors:
> 1. Push to GitHub (main/develop/dockerimage branch)
> 2. Wait for CI to complete
> 3. Make the package public (GitHub repo → Packages → Package settings → Change visibility)
> 4. OR build locally: `docker build -f Dockerfile.spaceapi -t spaceapi:local .`

#### Quick Start

1. **Create your `spaceapi.json` configuration file** (use `spaceapi.json.example` as template):
   ```bash
   curl -O https://raw.githubusercontent.com/q30-space/spaceapi-endpoint/main/spaceapi.json.example
   mv spaceapi.json.example spaceapi.json
   # Edit spaceapi.json with your space details
   ```

2. **Create `.env` file with your API key**:
   ```bash
   echo "SPACEAPI_AUTH_KEY=$(openssl rand -hex 32)" > .env
   ```

3. **Run the container**:
   ```bash
   docker run -d \
     --name spaceapi \
     -p 8080:8080 \
     -v $(pwd)/spaceapi.json:/app/spaceapi.json:ro \
     --env-file .env \
     --restart unless-stopped \
     ghcr.io/q30-space/spaceapi-endpoint:latest
   ```

4. **Test the deployment**:
   ```bash
   curl http://localhost:8080/api/space
   curl http://localhost:8080/health
   ```

#### Docker Compose Setup

The repository includes a production-ready docker-compose configuration (`docker-compose.prod.yml`):

```bash
# Use the provided production compose file
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Stop the service
docker-compose -f docker-compose.prod.yml down
```

Or create your own `docker-compose.yml`:

```yaml
version: '3.8'

services:
  spaceapi:
    image: ghcr.io/q30-space/spaceapi-endpoint:latest
    container_name: spaceapi
    ports:
      - "8080:8080"
    volumes:
      - ./spaceapi.json:/app/spaceapi.json:ro
    env_file:
      - .env
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
```

Then run:
```bash
docker-compose up -d
```

#### Available Image Tags

- `latest` - Latest stable release from main branch
- `main` - Latest commit from main branch
- `develop` - Latest commit from develop branch
- `v1.0.0` - Specific version tags (when available)
- `main-<sha>` - Specific commit SHA

#### Multi-Architecture Support

The Docker images support both `amd64` and `arm64` architectures, so you can run them on:
- x86_64 / AMD64 servers
- ARM64 servers (including Raspberry Pi 4/5, AWS Graviton, etc.)

#### Behind a Reverse Proxy (Caddy/Nginx)

Example Caddy configuration:
```caddyfile
your-domain.com {
    reverse_proxy localhost:8080
    encode gzip
}
```

Example Nginx configuration:
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Alternative Deployment Methods

#### Docker Development
```bash
# Start the service
make docker-compose-up
# or
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f spaceapi
```

#### Production - From sources
1. Update `spaceapi.json` with your space details
2. Build binaries: `make build`
3. Deploy the `bin/spaceapi` binary
4. Optionally deploy `bin/spaceicon` for monitoring

#### Production - Docker from source
1. Update `spaceapi.json` with your space details
2. Build binaries: `make build`
3. Start the container: `make docker-compose-up` (or just `docker compose up -d`)

#### Production - Caddy from source
1. Update `spaceapi.json` with your space details
2. Update `Caddyfile`
3. Copy `Caddyfile`, `docker-compose-caddy.yml` and `.env` to the parent directory and cd there
4. Build binaries: `docker compose build --no-cache -f docker-compose-caddy.yml spaceapi`
5. Start the containers: `docker compose up -d -f docker-compose-caddy.yml`

## Security Considerations

1. **API Authentication**: POST endpoints now require API key authentication via `SPACEAPI_AUTH_KEY` environment variable.
2. **Rate Limiting**: Failed authentication attempts are rate limited (5 attempts in 15 minutes = 1 hour block).
3. **HTTPS Required**: Production deployments must use HTTPS to protect API keys in transit.
4. **CORS**: Currently allows all origins. Restrict in production.
5. **Input Validation**: Basic validation is implemented, but consider additional checks.
6. **Key Management**: API keys should be rotated regularly and stored securely.


### Example JavaScript Integration

```javascript
// Fetch space status
fetch('http://localhost:8089/api/space')
  .then(response => response.json())
  .then(data => {
    const statusElement = document.getElementById('space-status');
    if (data.state && data.state.open) {
      statusElement.textContent = 'Space is OPEN';
      statusElement.className = 'status-open';
    } else {
      statusElement.textContent = 'Space is CLOSED';
      statusElement.className = 'status-closed';
    }
  });
```

## Monitoring

- **Health Checks**: Both services have health check endpoints
- **Logs**: Check Docker logs for issues
- **Metrics**: Consider adding Prometheus metrics for production

## Troubleshooting

### Common Issues

1. **API not responding**: Check if the spaceapi service is running
2. **CORS errors**: Ensure the API URL is correct
3. **JSON parsing errors**: Validate your JSON payloads

### Debug Commands

```bash
# Check service status
docker-compose ps

# View API logs
docker-compose logs spaceapi

# Test API endpoint
curl http://localhost:8089/health

# Test SpaceAPI endpoint
curl http://localhost:8089/api/space
```

## Project Structure

```
spaceapi-endpoint/
├── cmd/
│   ├── spaceapi/          # SpaceAPI server
│   └── spaceicon/         # CLI status tool
├── internal/
│   ├── handlers/          # HTTP handlers
│   │   └── spaceapi_test.go  # Handler tests
│   ├── middleware/        # Auth, CORS middleware
│   │   ├── auth_test.go      # Authentication tests
│   │   └── cors_test.go      # CORS tests
│   ├── models/           # Data models
│   ├── services/         # Business logic
│   └── testutil/         # Test helpers
│       └── helpers.go       # Mock data functions
├── scripts/              # Management scripts
├── bin/                  # Built binaries
├── spaceapi.json         # Configuration
└── Makefile             # Build automation
```

## Future Enhancements

1. **Webhooks**: Add webhook support for external integrations
2. **Metrics**: Add Prometheus metrics endpoint for monitoring
