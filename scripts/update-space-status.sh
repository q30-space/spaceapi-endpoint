#!/bin/bash

# SpaceAPI Status Update Script
# Usage: ./update-space-status.sh [open|closed] [message] [person]

SPACEAPI_URL="${SPACEAPI_URL:-https://example.com/spaceapi}"
SPACEAPI_AUTH_KEY="${SPACEAPI_AUTH_KEY:-}"
ACTION="${1:-open}"
MESSAGE="${2:-}"
PERSON="${3:-}"

# Check if API key is provided
if [ -z "$SPACEAPI_AUTH_KEY" ]; then
    echo "Error: SPACEAPI_AUTH_KEY environment variable is required"
    echo "Set it in your .env file or export it: export SPACEAPI_AUTH_KEY=your_key_here"
    exit 1
fi

case $ACTION in
    "open")
        OPEN=true
        DEFAULT_MESSAGE="Space is open"
        ;;
    "closed")
        OPEN=false
        DEFAULT_MESSAGE="Space is closed"
        ;;
    *)
        echo "Usage: $0 [open|closed] [message] [person]"
        echo "Example: $0 open 'Open for members' 'John Doe'"
        exit 1
        ;;
esac

# Use default message if none provided
if [ -z "$MESSAGE" ]; then
    MESSAGE="$DEFAULT_MESSAGE"
fi

# Prepare JSON payload
JSON_PAYLOAD=$(cat <<EOF
{
    "open": $OPEN,
    "message": "$MESSAGE",
    "trigger_person": "$PERSON"
}
EOF
)

# Send request
echo "Updating space status..."
echo "URL: $SPACEAPI_URL/api/space/state"
echo "Payload: $JSON_PAYLOAD"

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
    -d "$JSON_PAYLOAD" \
    "$SPACEAPI_URL/api/space/state")

# Extract HTTP status code and response body
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RESPONSE" | head -n -1)

case $HTTP_CODE in
    200)
        echo "✅ Space status updated successfully!"
        echo "Response: $RESPONSE_BODY"
        ;;
    401)
        echo "❌ Authentication failed: Invalid API key"
        echo "Check your SPACEAPI_AUTH_KEY environment variable"
        exit 1
        ;;
    403)
        echo "❌ Access forbidden: API key format invalid"
        exit 1
        ;;
    429)
        echo "❌ Rate limited: Too many failed authentication attempts"
        echo "Please wait before trying again"
        exit 1
        ;;
    *)
        echo "❌ Failed to update space status (HTTP $HTTP_CODE)"
        echo "Response: $RESPONSE_BODY"
        exit 1
        ;;
esac
