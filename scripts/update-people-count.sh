#!/bin/bash
# Copyright (C) 2025  pliski@q30.space
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

# SpaceAPI People Count Update Script
# Usage: ./update-people-count.sh [count] [location]

SPACEAPI_URL="${SPACEAPI_URL:-http://localhost:8089}"
SPACEAPI_AUTH_KEY="${SPACEAPI_AUTH_KEY:-}"
COUNT="${1:-0}"
LOCATION="${2:-Main Space}"

# Check if API key is provided
if [ -z "$SPACEAPI_AUTH_KEY" ]; then
    echo "Error: SPACEAPI_AUTH_KEY environment variable is required"
    echo "Set it in your .env file or export it: export SPACEAPI_AUTH_KEY=your_key_here"
    exit 1
fi

# Validate count is a number
if ! [[ "$COUNT" =~ ^[0-9]+$ ]]; then
    echo "Error: Count must be a number"
    echo "Usage: $0 [count] [location]"
    echo "Example: $0 5 'Main Space'"
    exit 1
fi

# Prepare JSON payload
JSON_PAYLOAD=$(cat <<EOF
{
    "value": $COUNT,
    "location": "$LOCATION"
}
EOF
)

# Send request
echo "Updating people count..."
echo "URL: $SPACEAPI_URL/api/space/people"
echo "Payload: $JSON_PAYLOAD"

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
    -d "$JSON_PAYLOAD" \
    "$SPACEAPI_URL/api/space/people")

# Extract HTTP status code and response body
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RESPONSE" | head -n -1)

case $HTTP_CODE in
    200)
        echo "✅ People count updated successfully!"
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
        echo "❌ Failed to update people count (HTTP $HTTP_CODE)"
        echo "Response: $RESPONSE_BODY"
        exit 1
        ;;
esac
