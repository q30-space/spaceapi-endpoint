# SpaceAPI Docker Deployment Guide

This guide will help you deploy the SpaceAPI server using the pre-built Docker image.

## Prerequisites

- Docker installed on your system
- (Optional) Docker Compose for easier management

## Important: First-Time Setup

### Option 1: Wait for CI to Build (Recommended for Production)

1. Copy [docker-compose-prod.yml](../docker-compose.prod.yml) to your host.
2. Create your `.env` and `spaceapi.json` files as described in [README](../README.md).
3. `docker-compose up -d`

**Making the Package Public:**
1. Go to your GitHub repository page
2. Click on "Packages" in the right sidebar
3. Click on the `spaceapi-endpoint` package
4. Click "Package settings"
5. Scroll to "Danger Zone" → "Change visibility"
6. Change to "Public"

### Option 2: Build Locally (Quick Testing)

If you want to test immediately before CI runs:

```bash
# Build the image locally
docker build -f Dockerfile.spaceapi -t spaceapi:local .

# Run it with your local tag
docker run -d \
  --name spaceapi \
  -p 8080:8080 \
  -v $(pwd)/spaceapi.json:/app/spaceapi.json:ro \
  --env-file .env \
  --restart unless-stopped \
  spaceapi:local
```

## Quick Start (After Image is Published)

### 1. Download the Example Configuration

```bash
# Create a directory for your deployment
mkdir spaceapi-server
cd spaceapi-server

# Download the example configuration
curl -O https://raw.githubusercontent.com/q30-space/spaceapi-endpoint/main/spaceapi.json.example

# Rename it to spaceapi.json
mv spaceapi.json.example spaceapi.json
```

### 2. Configure Your Space

Edit `spaceapi.json` and update at minimum:

- `space` - Your hackerspace name
- `url` - Your website URL
- `logo` - Your logo URL
- `location` - Your physical address and coordinates
- `contact` - Your contact information

### 3. Generate API Key

```bash
# Generate a secure API key
echo "SPACEAPI_AUTH_KEY=$(openssl rand -hex 32)" > .env

# View the generated key (save it somewhere safe!)
cat .env
```

### 4. Start the Server

**Option A: Using Docker directly**

```bash
docker run -d \
  --name spaceapi \
  -p 8080:8080 \
  -v $(pwd)/spaceapi.json:/app/spaceapi.json:ro \
  --env-file .env \
  --restart unless-stopped \
  ghcr.io/q30-space/spaceapi-endpoint:latest
```

**Option B: Using Docker Compose**

```bash
# Download the production compose file
curl -O https://raw.githubusercontent.com/q30-space/spaceapi-endpoint/main/docker-compose.prod.yml

# Start the service
docker-compose -f docker-compose.prod.yml up -d
```

### 5. Verify It Works

```bash
# Check the health endpoint
curl http://localhost:8080/health

# Check your SpaceAPI endpoint
curl http://localhost:8080/api/space
```

You should see your space information in JSON format!

## Updating the Space Status

Now you can update your space status using the API:

```bash
# Get your API key from the .env file
source .env

# Open the space
curl -X POST \
  -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
  -H "Content-Type: application/json" \
  -d '{"open": true, "message": "Space is open!"}' \
  http://localhost:8080/api/space/state

# Close the space
curl -X POST \
  -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
  -H "Content-Type: application/json" \
  -d '{"open": false, "message": "Space is closed"}' \
  http://localhost:8080/api/space/state

# Update people count
curl -X POST \
  -H "X-API-Key: $SPACEAPI_AUTH_KEY" \
  -H "Content-Type: application/json" \
  -d '{"value": 5, "location": "Main Space"}' \
  http://localhost:8080/api/space/people
```

## Production Deployment

### Behind a Reverse Proxy

For production, you should run the service behind a reverse proxy with HTTPS.

**Caddy (Recommended - automatic HTTPS)**

```caddyfile
api.yourdomain.com {
    reverse_proxy localhost:8080
    encode gzip
}
```

**Nginx with Let's Encrypt**

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Updating the Image

```bash
# Pull the latest image
docker pull ghcr.io/q30-space/spaceapi-endpoint:latest

# Recreate the container
docker-compose -f docker-compose.prod.yml up -d --force-recreate

# Or with plain docker
docker stop spaceapi
docker rm spaceapi
# Then run the docker run command again from Quick Start
```

### Viewing Logs

```bash
# With Docker Compose
docker-compose -f docker-compose.prod.yml logs -f

# With plain Docker
docker logs -f spaceapi
```

### Backup Your Configuration

Make sure to backup your files regularly:

- `spaceapi.json` - Your space configuration
- `.env` - Your API key

```bash
# Create a backup
tar czf spaceapi-backup-$(date +%Y%m%d).tar.gz spaceapi.json .env
```

## Troubleshooting

### Container won't start

```bash
# Check the logs
docker logs spaceapi

# Common issues:
# - Missing spaceapi.json file
# - Invalid JSON in spaceapi.json
# - Port 8080 already in use
```

### Can't reach the API

```bash
# Check if the container is running
docker ps | grep spaceapi

# Check if the port is accessible
curl http://localhost:8080/health

# If using a firewall, make sure port 8080 is open
sudo ufw allow 8080/tcp
```

### API key authentication failing

```bash
# Verify your API key is set
docker exec spaceapi env | grep SPACEAPI_AUTH_KEY

# Make sure you're using the correct header
# -H "X-API-Key: your_key_here"
```

## Advanced Configuration

### Custom Port

```bash
# Edit docker-compose.prod.yml and change the ports section:
ports:
  - "9000:8080"  # Host port:Container port

# Or with docker run:
docker run -d \
  --name spaceapi \
  -p 9000:8080 \
  -v $(pwd)/spaceapi.json:/app/spaceapi.json:ro \
  --env-file .env \
  --restart unless-stopped \
  ghcr.io/q30-space/spaceapi-endpoint:latest
```

### Resource Limits

Add to your docker-compose.yml:

```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 128M
    reservations:
      cpus: '0.1'
      memory: 64M
```

## Support

- **Documentation**: See [README.md](../README.md) for full documentation
- **Issues**: Report bugs on [GitHub Issues](https://github.com/q30-space/spaceapi-endpoint/issues)
- **SpaceAPI Spec**: https://spaceapi.io/docs/

## Security Best Practices

1. ✅ Always use HTTPS in production (via reverse proxy)
2. ✅ Keep your API key secret (never commit `.env` to git)
3. ✅ Use strong API keys (32+ random characters)
4. ✅ Mount spaceapi.json as read-only (`:ro` flag)
5. ✅ Keep the Docker image updated
6. ✅ Use firewall rules to restrict access if needed
7. ✅ Monitor logs for suspicious activity
8. ✅ Rotate API keys periodically

## Next Steps

- Set up automatic space status updates (see [scripts/](../scripts/))
- Configure monitoring (health checks, uptime monitoring)
- Join the SpaceAPI community
- Add your space to the SpaceAPI directory


## Misc


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

#### Release Information

The project follows [Semantic Versioning](https://semver.org/) for releases:
- **Major versions** (v2.0.0): Breaking changes
- **Minor versions** (v1.1.0): New features, backward compatible
- **Patch versions** (v1.0.1): Bug fixes, backward compatible

**Current Status**: Pre-release (v0.x.x) - API may change before v1.0.0

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
2. Build binary: `make build`
3. Deploy the `bin/spaceapi` binary

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
