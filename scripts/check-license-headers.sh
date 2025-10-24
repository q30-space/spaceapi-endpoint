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

# Script to check that all source files have the required license header
# Usage: ./scripts/check-license-headers.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# License header template (first 15 lines)
LICENSE_HEADER_LINES=15
EXPECTED_COPYRIGHT="Copyright (C) 2025  pliski@q30.space"
EXPECTED_GPL_LINK="https://www.gnu.org/licenses/"

# Files to check (source code files)
SOURCE_FILES=(
    "cmd/spaceapi/main.go"
    "cmd/spaceicon/main.go"
    "internal/handlers/spaceapi.go"
    "internal/handlers/spaceapi_test.go"
    "internal/middleware/auth.go"
    "internal/middleware/cors.go"
    "internal/middleware/cors_test.go"
    "internal/models/spaceapi.go"
    "internal/services/spaceapi.go"
    "internal/testutil/helpers.go"
    "scripts/update-space-status.sh"
    "scripts/update-people-count.sh"
    "test-docker-local.sh"
    "Makefile"
    "Caddyfile"
)

# Function to check if a file has the license header
check_license_header() {
    local file="$1"
    local errors=0
    
    if [[ ! -f "$file" ]]; then
        echo -e "${RED}❌ File not found: $file${NC}"
        return 1
    fi
    
    # Check if file has enough lines
    local line_count=$(wc -l < "$file")
    if [[ $line_count -lt $LICENSE_HEADER_LINES ]]; then
        echo -e "${RED}❌ $file: File too short (${line_count} lines, need at least ${LICENSE_HEADER_LINES})${NC}"
        return 1
    fi
    
    # Check for copyright notice
    if ! head -n $LICENSE_HEADER_LINES "$file" | grep -q "$EXPECTED_COPYRIGHT"; then
        echo -e "${RED}❌ $file: Missing or incorrect copyright notice${NC}"
        errors=$((errors + 1))
    fi
    
    # Check for GPL license link
    if ! head -n $LICENSE_HEADER_LINES "$file" | grep -q "$EXPECTED_GPL_LINK"; then
        echo -e "${RED}❌ $file: Missing or incorrect GPL license link${NC}"
        errors=$((errors + 1))
    fi
    
    # Check for proper comment style based on file type
    local extension="${file##*.}"
    case "$extension" in
        "go")
            if ! head -n 1 "$file" | grep -q "^// Copyright"; then
                echo -e "${RED}❌ $file: Go files should start with // comment style${NC}"
                errors=$((errors + 1))
            fi
            ;;
        "sh")
            if ! head -n 2 "$file" | grep -q "^# Copyright"; then
                echo -e "${RED}❌ $file: Shell scripts should start with # comment style${NC}"
                errors=$((errors + 1))
            fi
            ;;
        *)
            # For Makefile, Caddyfile, and other files
            if ! head -n 1 "$file" | grep -q "^# Copyright"; then
                echo -e "${RED}❌ $file: Should start with # comment style${NC}"
                errors=$((errors + 1))
            fi
            ;;
    esac
    
    if [[ $errors -eq 0 ]]; then
        echo -e "${GREEN}✅ $file: License header OK${NC}"
        return 0
    else
        return 1
    fi
}

# Main execution
echo -e "${YELLOW}🔍 Checking license headers in source files...${NC}"
echo

total_files=0
failed_files=0

for file in "${SOURCE_FILES[@]}"; do
    total_files=$((total_files + 1))
    if ! check_license_header "$file"; then
        failed_files=$((failed_files + 1))
    fi
done

echo
echo "=========================================="
echo -e "📊 Summary: ${total_files} files checked, ${failed_files} failed"

if [[ $failed_files -eq 0 ]]; then
    echo -e "${GREEN}🎉 All source files have proper license headers!${NC}"
    exit 0
else
    echo -e "${RED}💥 $failed_files file(s) are missing or have incorrect license headers${NC}"
    echo
    echo "To fix this, ensure each source file starts with:"
    echo "  Copyright (C) 2025  pliski@q30.space"
    echo "  [GPL v3 license text]"
    echo "  https://www.gnu.org/licenses/"
    echo
    echo "Use the appropriate comment style:"
    echo "  - Go files: // comments"
    echo "  - Shell scripts: # comments"
    echo "  - Makefile/Caddyfile: # comments"
    exit 1
fi
