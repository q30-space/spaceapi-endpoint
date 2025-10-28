# OpenAPI Documentation

This document explains how to use the OpenAPI specification for the SpaceAPI Endpoint.

## Overview

The SpaceAPI Endpoint API is fully documented using the [OpenAPI Specification (OAS) 3.0](https://swagger.io/specification/) standard. The specification file is located at [`openapi.yaml`](../openapi.yaml) in the project root.

## Viewing the API Documentation

### Option 1: Swagger UI (Recommended)

The easiest way to view and interact with the API documentation is using Swagger UI:

```bash
# Start Swagger UI (requires Docker)
make openapi-ui
```

This will start Swagger UI on http://localhost:8081 where you can:
- Browse all API endpoints
- See request/response schemas
- Try out API calls directly
- View examples

### Option 2: Swagger Editor Online

1. Go to https://editor.swagger.io/
2. Click **File** → **Import file**
3. Upload the `openapi.yaml` file
4. The editor will display the documentation and allow you to make changes

### Option 3: Redoc

For a more polished, read-only documentation view:

```bash
# Using Docker
docker run -p 8082:80 \
  -e SPEC_URL=https://raw.githubusercontent.com/q30-space/spaceapi-endpoint/main/openapi.yaml \
  redocly/redoc
```

Then visit http://localhost:8082

### Option 4: VS Code Extension

Install the [OpenAPI (Swagger) Editor](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi) extension:

1. Open VS Code
2. Install "OpenAPI (Swagger) Editor" by 42Crunch
3. Open `openapi.yaml`
4. Click the preview icon in the top-right corner

## Validating the OpenAPI Specification

To ensure the OpenAPI specification is valid:

```bash
# Using the Makefile (requires Docker)
make openapi-validate
```

This uses the OpenAPI Generator CLI Docker image which is publicly available and doesn't require any authentication or npm dependencies.

## Generating Client Libraries

You can generate client libraries in various programming languages using the OpenAPI Generator:

### Install OpenAPI Generator

```bash
# Using npm
npm install -g @openapitools/openapi-generator-cli

# Using Docker (no installation needed)
# See examples below
```

### Generate Client Libraries

#### JavaScript/TypeScript Client

```bash
docker run --rm \
  -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate \
  -i /local/openapi.yaml \
  -g typescript-fetch \
  -o /local/clients/typescript
```

#### Python Client

```bash
docker run --rm \
  -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate \
  -i /local/openapi.yaml \
  -g python \
  -o /local/clients/python
```

#### Go Client

```bash
docker run --rm \
  -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate \
  -i /local/openapi.yaml \
  -g go \
  -o /local/clients/go
```

#### Java Client

```bash
docker run --rm \
  -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate \
  -i /local/openapi.yaml \
  -g java \
  -o /local/clients/java
```

### Available Generators

View all available client generators:
```bash
docker run --rm openapitools/openapi-generator-cli list
```

## Generating Server Stubs

You can also generate server stubs in various frameworks:

```bash
# Node.js/Express server
docker run --rm \
  -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate \
  -i /local/openapi.yaml \
  -g nodejs-express-server \
  -o /local/server-stubs/nodejs

# Spring Boot server
docker run --rm \
  -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate \
  -i /local/openapi.yaml \
  -g spring \
  -o /local/server-stubs/spring
```

## Generating HTML Documentation

To generate a standalone HTML documentation file:

```bash
# Using redoc-cli
npm install -g redoc-cli
redoc-cli bundle openapi.yaml -o api-documentation.html

# Or using Docker
docker run --rm \
  -v ${PWD}:/spec \
  redocly/redoc-cli bundle /spec/openapi.yaml \
  -o /spec/api-documentation.html
```

## API Endpoints Summary

### Public Endpoints (No Authentication)

- `GET /` - Get SpaceAPI data (root endpoint)
- `GET /api/space` - Get SpaceAPI data
- `GET /health` - Health check

### Protected Endpoints (Require API Key)

All POST endpoints require authentication via the `X-API-Key` header:

- `POST /api/space/state` - Update space open/closed state
- `POST /api/space/people` - Update people count
- `POST /api/space/event` - Add an event

## Authentication

Protected endpoints require an API key passed in the `X-API-Key` header:

```bash
curl -X POST \
  -H "X-API-Key: your_api_key_here" \
  -H "Content-Type: application/json" \
  -d '{"open": true}' \
  http://localhost:8080/api/space/state
```

Set your API key in the `.env` file:
```
SPACEAPI_AUTH_KEY=your_generated_key_here
```

## Rate Limiting

Failed authentication attempts are rate limited:
- **Limit**: 5 failed attempts within 15 minutes
- **Block Duration**: 1 hour
- **Response**: HTTP 429 with `Retry-After` header
- **Scope**: Per IP address

## Error Responses

| Status | Description | Response Body |
|--------|-------------|---------------|
| 400 | Invalid JSON | `"Invalid JSON"` |
| 401 | Missing/invalid API key | `"Invalid API key"` |
| 403 | Access forbidden | `"Access forbidden"` |
| 429 | Rate limited | `"Too many failed authentication attempts..."` |

## Integration Examples

### cURL

```bash
# Get space status
curl http://localhost:8080/api/space

# Update state (requires auth)
curl -X POST \
  -H "X-API-Key: your_key" \
  -H "Content-Type: application/json" \
  -d '{"open": true, "message": "Space is open"}' \
  http://localhost:8080/api/space/state
```

### JavaScript (Fetch API)

```javascript
// Get space status
fetch('http://localhost:8080/api/space')
  .then(response => response.json())
  .then(data => console.log(data));

// Update state (requires auth)
fetch('http://localhost:8080/api/space/state', {
  method: 'POST',
  headers: {
    'X-API-Key': 'your_key',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    open: true,
    message: 'Space is open'
  })
})
  .then(response => response.json())
  .then(data => console.log(data));
```

### Python (requests)

```python
import requests

# Get space status
response = requests.get('http://localhost:8080/api/space')
data = response.json()
print(data)

# Update state (requires auth)
headers = {
    'X-API-Key': 'your_key',
    'Content-Type': 'application/json'
}
payload = {
    'open': True,
    'message': 'Space is open'
}
response = requests.post(
    'http://localhost:8080/api/space/state',
    headers=headers,
    json=payload
)
print(response.json())
```

## Updating the OpenAPI Specification

When you modify the API:

1. Update `openapi.yaml` to reflect the changes
2. Validate the specification: `make openapi-validate`
3. Test the changes with Swagger UI: `make openapi-ui`
4. Update this documentation if needed
5. Regenerate client libraries if you maintain them

## Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger Editor](https://editor.swagger.io/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Redoc](https://redocly.com/redoc/)
- [OpenAPI Generator](https://openapi-generator.tech/)
- [Spectral (OpenAPI Linter)](https://stoplight.io/open-source/spectral)
- [SpaceAPI v15 Specification](https://spaceapi.io/docs/)

## Contributing

When contributing to the API:

1. Ensure all changes are reflected in `openapi.yaml`
2. Validate the OpenAPI spec before committing
3. Update examples in the specification
4. Update this documentation
5. Test generated client libraries if applicable

