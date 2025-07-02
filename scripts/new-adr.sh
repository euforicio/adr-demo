#!/bin/bash

# Create New ADR
# This script creates a new ADR from the template with the next sequential number

set -e

# Change to repository root
cd "$(dirname "$0")/.."

# Check if template exists
if [ ! -f "docs/adr/template.md" ]; then
  echo "❌ Template file not found: docs/adr/template.md"
  exit 1
fi

# Find the next ADR number
last_number=$(find docs/adr -name "[0-9][0-9][0-9][0-9]-*.md" -exec basename {} \; | sed 's/^\([0-9]\{4\}\)-.*/\1/' | sort -n | tail -1)

if [ -z "$last_number" ]; then
  next_number="0001"
else
  next_number=$(printf "%04d" $((10#$last_number + 1)))
fi

echo "📄 Creating new ADR $next_number"

# Get title from user
if [ $# -eq 0 ]; then
  echo "Please enter the ADR title (use imperative mood, e.g., 'Use Redis for caching'):"
  read -r title
else
  title="$*"
fi

# Convert title to filename format (kebab-case)
filename_title=$(echo "$title" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-\|-$//g')

# Create filename
filename="docs/adr/${next_number}-${filename_title}.md"

# Check if file already exists
if [ -f "$filename" ]; then
  echo "❌ File already exists: $filename"
  exit 1
fi

# Copy template and customize
cp "docs/adr/template.md" "$filename"

# Replace placeholders in the new file
sed -i.bak "s/# \[Title\]/# $title/" "$filename"
sed -i.bak "s/\[Brief description of the decision\]/Brief description of the decision to $filename_title/" "$filename"
sed -i.bak "s/\[Current date\]/$(date '+%Y-%m-%d')/" "$filename"

# Clean up backup file
rm -f "${filename}.bak"

echo "✅ Created new ADR: $filename"
echo ""
echo "📝 Next steps:"
echo "   1. Edit the ADR content: $filename"
echo "   2. Update the status from 'Proposed' when decided"
echo "   3. Add to git: git add $filename"
echo "   4. Update README index: make generate-index"
echo ""
echo "💡 To open in editor:"
echo "   \$EDITOR $filename"

# Optionally open in editor
if [ -n "$EDITOR" ]; then
  echo ""
  read -p "Open in \$EDITOR now? (y/N): " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    $EDITOR "$filename"
  fi
fi