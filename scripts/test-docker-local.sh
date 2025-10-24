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
set -e

echo "🐳 Building Docker image locally..."
docker build -f Dockerfile.spaceapi -t spaceapi:local .

echo ""
echo "✅ Docker image built successfully!"
echo ""
echo "📝 To run it, make sure you have:"
echo "   1. A spaceapi.json file in the current directory"
echo "   2. A .env file with SPACEAPI_AUTH_KEY"
echo ""
echo "Then run:"
echo ""
echo "  docker run -d \\"
echo "    --name spaceapi \\"
echo "    -p 8080:8080 \\"
echo "    -v \$(pwd)/spaceapi.json:/app/spaceapi.json:ro \\"
echo "    --env-file .env \\"
echo "    --restart unless-stopped \\"
echo "    spaceapi:local"
echo ""
echo "Or use the example config for testing:"
echo ""
echo "  cp spaceapi.json.example spaceapi.json"
echo "  echo 'SPACEAPI_AUTH_KEY=test-key-12345' > .env"
echo "  docker run -d \\"
echo "    --name spaceapi \\"
echo "    -p 8080:8080 \\"
echo "    -v \$(pwd)/spaceapi.json:/app/spaceapi.json:ro \\"
echo "    --env-file .env \\"
echo "    spaceapi:local"
