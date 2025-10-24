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
# Usage: 
#   ./scripts/check-license-headers.sh         # Check mode (for CI)
#   ./scripts/check-license-headers.sh --fix   # Fix mode (adds missing headers)

set -e

# Parse command line arguments
FIX_MODE=false
if [[ "$1" == "--fix" ]]; then
    FIX_MODE=true
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# License header template (first 15 lines)
LICENSE_HEADER_LINES=15
EXPECTED_COPYRIGHT="Copyright (C) 2025  pliski@q30.space"
EXPECTED_GPL_LINK="https://www.gnu.org/licenses/"

# Get the project root directory (one level up from scripts/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Additional directories to exclude beyond .gitignore
EXTRA_EXCLUDE_DIRS=("doc")

# Function to find all source files dynamically
# Uses git ls-files to respect .gitignore automatically
find_source_files() {
    cd "$PROJECT_ROOT"
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${RED}Error: Not in a git repository. Cannot use git ls-files.${NC}" >&2
        exit 1
    fi
    
    # Use git ls-files to get all tracked files (respects .gitignore)
    # Then filter for source files we care about
    git ls-files | grep -E '\.(go|sh)$|^Makefile$|^Caddyfile$' | while read -r file; do
        # Apply extra exclusions
        local skip=false
        for exclude_dir in "${EXTRA_EXCLUDE_DIRS[@]}"; do
            if [[ "$file" == "$exclude_dir/"* ]]; then
                skip=true
                break
            fi
        done
        
        if [[ "$skip" == false ]]; then
            echo "$file"
        fi
    done
}

# Build array of source files to check
mapfile -t SOURCE_FILES < <(find_source_files | sort)

# License header templates
get_license_header() {
    local file_type="$1"
    local comment_prefix="$2"
    
    cat << EOF
${comment_prefix} Copyright (C) 2025  pliski@q30.space
${comment_prefix}
${comment_prefix} This program is free software: you can redistribute it and/or modify
${comment_prefix} it under the terms of the GNU General Public License as published by
${comment_prefix} the Free Software Foundation, either version 3 of the License, or
${comment_prefix} (at your option) any later version.
${comment_prefix}
${comment_prefix} This program is distributed in the hope that it will be useful,
${comment_prefix} but WITHOUT ANY WARRANTY; without even the implied warranty of
${comment_prefix} MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
${comment_prefix} GNU General Public License for more details.
${comment_prefix}
${comment_prefix} You should have received a copy of the GNU General Public License
${comment_prefix} along with this program.  If not, see <https://www.gnu.org/licenses/>.

EOF
}

# Function to add license header to a file
add_license_header() {
    local file="$1"
    local extension="${file##*.}"
    local comment_prefix
    local temp_file="${file}.tmp"
    
    # Determine comment style
    case "$extension" in
        "go")
            comment_prefix="//"
            ;;
        "sh")
            comment_prefix="#"
            # Preserve shebang if it exists
            if head -n 1 "$file" | grep -q "^#!"; then
                local shebang=$(head -n 1 "$file")
                local content=$(tail -n +2 "$file")
                {
                    echo "$shebang"
                    get_license_header "$extension" "$comment_prefix"
                    echo "$content"
                } > "$temp_file"
                mv "$temp_file" "$file"
                return 0
            fi
            ;;
        *)
            # For Makefile, Caddyfile, and other files
            comment_prefix="#"
            ;;
    esac
    
    # Add header to the beginning of the file
    {
        get_license_header "$extension" "$comment_prefix"
        cat "$file"
    } > "$temp_file"
    mv "$temp_file" "$file"
}

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
if [[ "$FIX_MODE" == true ]]; then
    echo -e "${YELLOW}🔧 Fixing license headers in source files...${NC}"
else
    echo -e "${YELLOW}🔍 Checking license headers in source files...${NC}"
fi
echo

total_files=0
failed_files=0
fixed_files=0

for file in "${SOURCE_FILES[@]}"; do
    total_files=$((total_files + 1))
    
    if ! check_license_header "$file" 2>/dev/null; then
        failed_files=$((failed_files + 1))
        
        if [[ "$FIX_MODE" == true ]]; then
            echo -e "${YELLOW}🔧 Adding license header to: $file${NC}"
            add_license_header "$file"
            fixed_files=$((fixed_files + 1))
            
            # Verify the fix worked
            if check_license_header "$file" 2>/dev/null; then
                echo -e "${GREEN}✅ $file: License header added successfully${NC}"
            else
                echo -e "${RED}❌ $file: Failed to add license header correctly${NC}"
            fi
        fi
    fi
done

echo
echo "=========================================="

if [[ "$FIX_MODE" == true ]]; then
    echo -e "📊 Summary: ${total_files} files checked, ${fixed_files} files fixed"
    if [[ $fixed_files -gt 0 ]]; then
        echo -e "${GREEN}🎉 Successfully added license headers to ${fixed_files} file(s)!${NC}"
        echo -e "${YELLOW}⚠️  Don't forget to review and commit the changes.${NC}"
    else
        echo -e "${GREEN}🎉 All source files already have proper license headers!${NC}"
    fi
    exit 0
else
    echo -e "📊 Summary: ${total_files} files checked, ${failed_files} failed"
    
    if [[ $failed_files -eq 0 ]]; then
        echo -e "${GREEN}🎉 All source files have proper license headers!${NC}"
        exit 0
    else
        echo -e "${RED}💥 $failed_files file(s) are missing or have incorrect license headers${NC}"
        echo
        echo "To automatically fix these files, run:"
        echo -e "  ${YELLOW}./scripts/check-license-headers.sh --fix${NC}"
        echo
        echo "Or manually ensure each source file starts with:"
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
fi
