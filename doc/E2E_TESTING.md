# End-to-End Testing

## Overview

The CI pipeline includes comprehensive end-to-end (E2E) tests that validate the Docker container against the OpenAPI specification. These tests ensure that the deployed service correctly implements the API contract.

## Test Framework

The E2E tests use [Schemathesis](https://schemathesis.readthedocs.io/), a property-based testing tool specifically designed for testing APIs against OpenAPI/Swagger specifications.

### Why Schemathesis?

- **Automatic test generation**: Generates test cases from the OpenAPI spec
- **Property-based testing**: Uses Hypothesis to generate edge cases automatically
- **OpenAPI-native**: Understands OpenAPI schemas and validates responses
- **Comprehensive validation**: Checks response schemas, status codes, headers, etc.

## Test Coverage

The E2E test job (`e2e-test`) includes:

### 1. OpenAPI Conformance Tests (Public Endpoints)

Tests public endpoints without authentication:
- `GET /health` - Health check endpoint
- `GET /` - Root SpaceAPI endpoint
- `GET /api/space` - SpaceAPI data endpoint

Validates:
- Response schemas match OpenAPI spec
- Status codes are correct
- Content types are correct
- Response structure is valid

### 2. OpenAPI Conformance Tests (Protected Endpoints)

Tests protected endpoints with valid API key:
- `POST /api/space/state` - Update space state
- `POST /api/space/people` - Update people count
- `POST /api/space/event` - Add event

Validates:
- Authentication works correctly
- Request/response schemas match spec
- All documented fields are handled properly

### 3. Authentication Tests

Explicit tests for authentication error cases:
- Missing API key (expects 401)
- Invalid API key (expects 401)

Ensures security is properly implemented.

### 4. Functional API Operation Tests

Tests the complete flow of API operations:
- Update space state to open
- Update people count
- Create check-in event
- Verify all updates are reflected in GET responses

Validates end-to-end data persistence and retrieval.

## Running Tests Locally

### Prerequisites

```bash
# Install Schemathesis
pip install schemathesis requests

# Ensure you have Docker and docker-compose installed
```

### Run the full E2E test suite

```bash
# 1. Prepare test configuration
cp spaceapi.json.example spaceapi.json

# 2. Start the service
export SPACEAPI_AUTH_KEY="test-api-key"
docker-compose up -d

# 3. Wait for service to be healthy
until docker-compose ps | grep -q "healthy"; do sleep 2; done

# 4. Run OpenAPI conformance tests (public endpoints)
schemathesis run openapi.yaml \
  --base-url http://localhost:8089 \
  --checks all \
  --hypothesis-phases=explicit \
  --workers 4 \
  --include-path-regex "^/(health|api/space)$" \
  --exclude-path-regex ".*/(state|people|event)$" \
  --show-errors-tracebacks \
  --report

# 5. Run OpenAPI conformance tests (protected endpoints)
schemathesis run openapi.yaml \
  --base-url http://localhost:8089 \
  --header "X-API-Key: test-api-key" \
  --checks all \
  --hypothesis-phases=explicit \
  --workers 4 \
  --include-path-regex ".*/(state|people|event)$" \
  --show-errors-tracebacks \
  --report

# 6. Test authentication (manual curl tests)
# Missing API key (should return 401)
curl -X POST http://localhost:8089/api/space/state \
  -H "Content-Type: application/json" \
  -d '{"open": true}'

# Invalid API key (should return 401)
curl -X POST http://localhost:8089/api/space/state \
  -H "Content-Type: application/json" \
  -H "X-API-Key: invalid-key" \
  -d '{"open": true}'

# 7. Test successful operations
# Update state
curl -X POST http://localhost:8089/api/space/state \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test-api-key" \
  -d '{"open": true, "message": "Test", "trigger_person": "Tester"}'

# Update people count
curl -X POST http://localhost:8089/api/space/people \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test-api-key" \
  -d '{"value": 5, "location": "Main Space"}'

# Add event
curl -X POST http://localhost:8089/api/space/event \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test-api-key" \
  -d '{"name": "Tester", "type": "check-in", "extra": "Testing"}'

# Verify updates
curl http://localhost:8089/api/space | jq

# 8. Cleanup
docker-compose down -v
```

## CI Pipeline Integration

The `e2e-test` job runs:
- **In parallel with** `docker-build` (after basic tests pass)
- **Before deployment** to catch issues early
- **On all branches** to validate PRs

### Job Dependencies

```
test, build, lint, license-check, openapi-validate
  ↓
e2e-test (runs in parallel with docker-build)
```

## Troubleshooting

### Tests fail with connection errors

Check if the service started correctly:
```bash
docker-compose logs
docker-compose ps
```

### Tests fail on specific endpoints

1. Check OpenAPI spec matches implementation
2. Review the Schemathesis output for details
3. Test endpoint manually with curl
4. Check authentication header is passed correctly

### Service doesn't become healthy

The healthcheck waits up to 60 seconds. If it fails:
1. Check `spaceapi.json` configuration is valid
2. Check Docker logs for startup errors
3. Verify the `/health` endpoint is accessible

## Extending Tests

### Add custom test scenarios

Create a new step in the `e2e-test` job:

```yaml
- name: Test custom scenario
  run: |
    # Your custom test logic here
    response=$(curl -s http://localhost:8089/your-endpoint)
    # Validate response
```

### Test additional OpenAPI properties

Modify Schemathesis checks:

```bash
schemathesis run openapi.yaml \
  --base-url http://localhost:8089 \
  --checks all \
  --hypothesis-max-examples 50 \  # More test cases
  --hypothesis-seed 1234 \         # Reproducible tests
  --stateful=links                 # Stateful testing
```

### Add load testing

Install and use `locust` or `k6` for load testing:

```yaml
- name: Load test
  run: |
    pip install locust
    locust -f load_test.py --headless -u 10 -r 1 -t 30s
```

## References

- [Schemathesis Documentation](https://schemathesis.readthedocs.io/)
- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

