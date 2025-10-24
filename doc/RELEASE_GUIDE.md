# Release Guide

This guide explains how to create releases for the SpaceAPI Endpoint project.

## Quick Start

### 1. Create Your First Release

```bash
# Make sure you're on the main branch and everything is committed
git checkout main
git pull origin main

# Create the first release (v0.1.0 for initial release)
./scripts/create-release.sh v0.1.0
```

### 2. What Happens Automatically

When you run the release script, GitHub Actions will:

1. **Build binaries** for multiple platforms:
   - Linux (AMD64, ARM64)
   - Windows (AMD64)
   - macOS (AMD64, ARM64)

2. **Create a GitHub release** with:
   - Download links for all binaries
   - Release notes
   - Checksums for verification

3. **Build and push Docker images** to GitHub Container Registry:
   - `ghcr.io/q30-space/spaceapi-endpoint:v0.1.0`
   - `ghcr.io/q30-space/spaceapi-endpoint:latest`
   - `ghcr.io/q30-space/spaceapi-endpoint:main`

### 3. Make Docker Image Public

After the release is created:

1. Go to [GitHub Packages](https://github.com/q30-space/spaceapi-endpoint/pkgs/container/spaceapi-endpoint)
2. Click "Package settings"
3. Scroll to "Danger Zone"
4. Click "Change visibility" → "Public"

## Release Process Details

### Prerequisites

- You must be on the `main` branch
- Working directory must be clean (no uncommitted changes)
- You need push access to the repository

### Version Format

Use [Semantic Versioning](https://semver.org/):
- `v0.1.0` - First release (pre-release)
- `v1.0.0` - First stable release
- `v1.1.0` - New features
- `v1.0.1` - Bug fixes

### Release Script Options

```bash
# Interactive mode (script will prompt for version)
./scripts/create-release.sh

# Specify version directly
./scripts/create-release.sh v1.0.0

# Using make
make release VERSION=v1.0.0
```

### Manual Release Process

If you need to create a release manually:

```bash
# 1. Update version in go.mod
go mod edit -module=github.com/q30-space/spaceapi-endpoint@1.0.0

# 2. Run tests and checks
make test
make lint
make check-license

# 3. Build binaries
make build

# 4. Commit and tag
git add .
git commit -m "chore: bump version to v1.0.0"
git tag -a v1.0.0 -m "Release v1.0.0"

# 5. Push
git push origin main
git push origin v1.0.0
```

## Docker Images

### Available Tags

- `latest` - Latest stable release
- `main` - Latest commit from main branch
- `develop` - Latest commit from develop branch
- `v1.0.0` - Specific version
- `main-abc1234` - Specific commit

### Using Docker Images

```bash
# Pull latest release
docker pull ghcr.io/q30-space/spaceapi-endpoint:latest

# Pull specific version
docker pull ghcr.io/q30-space/spaceapi-endpoint:v1.0.0

# Run container
docker run -p 8080:8080 \
  -v $(pwd)/spaceapi.json:/app/spaceapi.json:ro \
  --env-file .env \
  ghcr.io/q30-space/spaceapi-endpoint:latest
```

## Troubleshooting

### Common Issues

1. **"Tag already exists"**
   - The version tag already exists
   - Choose a different version or delete the existing tag

2. **"Working directory not clean"**
   - Commit or stash your changes first
   - `git status` to see what's uncommitted

3. **"Not on main branch"**
   - Switch to main branch: `git checkout main`
   - Or use `--force` flag (not recommended)

4. **Docker image not found**
   - Wait for GitHub Actions to complete (5-10 minutes)
   - Check if the package is public
   - Verify the tag was pushed correctly

### Checking Release Status

```bash
# Check if tag exists
git tag -l | grep v1.0.0

# Check GitHub Actions status
# Go to: https://github.com/q30-space/spaceapi-endpoint/actions

# Check Docker image
docker pull ghcr.io/q30-space/spaceapi-endpoint:v1.0.0
```

## Next Steps After First Release

1. **Test the release**:
   ```bash
   # Download and test binaries
   curl -L https://github.com/q30-space/spaceapi-endpoint/releases/download/v0.1.0/spaceapi-linux-amd64 -o spaceapi
   chmod +x spaceapi
   ./spaceapi --version
   ```

2. **Test Docker image**:
   ```bash
   docker run --rm ghcr.io/q30-space/spaceapi-endpoint:v0.1.0 --version
   ```

3. **Update documentation** if needed

4. **Announce the release** to your community

## Support

If you encounter issues:

1. Check the [GitHub Actions logs](https://github.com/q30-space/spaceapi-endpoint/actions)
2. Verify all prerequisites are met
3. Check the [troubleshooting section](#troubleshooting)
4. Open an issue on GitHub if needed


## TLDR: Creating Releases

### Automatic Release Process

The project uses GitHub Actions for automated releases:

1. **Create a version tag**:
   ```bash
   # Using the release script (recommended)
   ./scripts/create-release.sh v1.0.0
   
   # Or using make
   make release VERSION=v1.0.0
   
   # Or manually
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically**:
   - Build binaries for multiple platforms (Linux, Windows, macOS)
   - Create a GitHub release with download links
   - Build and push Docker images to GitHub Container Registry
   - Generate checksums for all binaries

### Manual Release Process

If you need to create a release manually:

1. **Ensure go.mod is up to date**:
   ```bash
   go mod tidy
   ```

2. **Run tests and checks**:
   ```bash
   make test
   make lint
   make check-license
   ```

3. **Build binaries**:
   ```bash
   make build
   ```

4. **Create and push tag**:
   ```bash
   git add .
   git commit -m "chore: bump version to v1.0.0"
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin main
   git push origin v1.0.0
   ```

### Docker Image Publishing

Docker images are automatically published to GitHub Container Registry:

- **Repository**: `ghcr.io/q30-space/spaceapi-endpoint`
- **Tags**: Version tags, `latest`, branch names
- **Architectures**: `linux/amd64`, `linux/arm64`
