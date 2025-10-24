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

# Script to create a new release
# Usage: ./scripts/create-release.sh [version]
# Example: ./scripts/create-release.sh v1.0.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository"
    exit 1
fi

# Check if we're on main branch
current_branch=$(git branch --show-current)
if [ "$current_branch" != "main" ]; then
    print_warning "Not on main branch (current: $current_branch)"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if working directory is clean
if ! git diff-index --quiet HEAD --; then
    print_error "Working directory is not clean. Please commit or stash your changes."
    exit 1
fi

# Get version from argument or prompt
if [ -n "$1" ]; then
    VERSION="$1"
else
    # Get the latest tag
    latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [ -n "$latest_tag" ]; then
        print_status "Latest tag: $latest_tag"
    else
        print_status "No previous tags found"
    fi
    
    read -p "Enter version (e.g., v1.0.0): " VERSION
fi

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Invalid version format. Expected format: v1.0.0"
    exit 1
fi

# Remove 'v' prefix for semantic versioning
SEMVER=${VERSION#v}

print_status "Creating release $VERSION (semver: $SEMVER)"

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    print_error "Tag $VERSION already exists"
    exit 1
fi

# Update version in go.mod if needed (optional)
print_status "Updating go.mod version..."
# Note: We don't actually need to change the module path for releases
# The module path should remain the same, only the version tag changes

# Run tests
print_status "Running tests..."
if ! make test; then
    print_error "Tests failed. Aborting release."
    exit 1
fi

# Run linting
print_status "Running linter..."
if ! make lint; then
    print_error "Linting failed. Aborting release."
    exit 1
fi

# Check license headers
print_status "Checking license headers..."
if ! make check-license; then
    print_error "License check failed. Aborting release."
    exit 1
fi

# Build binaries locally for testing
print_status "Building binaries..."
make build

# Test the binaries
print_status "Testing binaries..."
if ! ./bin/spaceapi --version; then
    print_error "spaceapi binary test failed"
    exit 1
fi

# Commit changes
print_status "Committing version update..."
git add go.mod
git commit -m "chore: bump version to $VERSION" || print_warning "No changes to commit"

# Create and push tag
print_status "Creating tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"

print_status "Pushing changes and tag..."
git push origin main
git push origin "$VERSION"

print_success "Release $VERSION created and pushed!"
print_status "GitHub Actions will now build and publish the release automatically."
print_status "Docker image will be available at: ghcr.io/q30-space/spaceapi-endpoint:$SEMVER"
print_status "Release page: https://github.com/q30-space/spaceapi-endpoint/releases/tag/$VERSION"

# Optional: Open release page
if command -v xdg-open > /dev/null; then
    read -p "Open release page in browser? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        xdg-open "https://github.com/q30-space/spaceapi-endpoint/releases/tag/$VERSION"
    fi
fi

print_success "Release process completed!"
