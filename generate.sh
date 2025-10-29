#!/bin/bash

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change to SDK root
cd "$SCRIPT_DIR" || exit 1

echo "📝 Generating content-sdk-go templ files..."

# Check if templ is installed
if ! command -v templ &> /dev/null; then
    echo "❌ templ CLI not found. Installing..."
    go install github.com/a-h/templ/cmd/templ@latest
fi

# Generate templ components
templ generate

echo "✅ Templ generation complete!"

