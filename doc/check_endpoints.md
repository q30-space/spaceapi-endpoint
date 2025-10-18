# Commands to check endpoints

**Note**: The POST endpoints require API key authentication. Set the `SPACEAPI_AUTH_KEY` environment variable before running the commands:
```sh
export SPACEAPI_AUTH_KEY=your_api_key_here
```

## Test SpaceAPI endpoint
```sh
curl -v http://localhost:8089/api/space
```

## Update space state (open)
```sh
curl -v -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
  -d '{"open": true, "message": "Space is open for members", "trigger_person": "John Doe"}' \
  http://localhost:8089/api/space/state
```

## Update space state (closed)
```sh
curl -v -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
  -d '{"open": false, "message": "Space is closed", "trigger_person": "Jane Smith"}' \
  http://localhost:8089/api/space/state
```

## Update people count
```sh
curl -v -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
  -d '{"value": 5, "location": "Main Space"}' \
  http://localhost:8089/api/space/people
```

## Add an event
```sh
curl -v -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
  -d '{"name": "John Doe", "type": "check-in", "extra": "Working on Arduino project"}' \
  http://localhost:8089/api/space/event
```

## Health check
```sh
curl -v http://localhost:8089/health
```