#!/bin/bash

# Generate ADR Index JSON
# This script scans the docs/adr directory and creates an adr-index.json file
# that the web application can use to load ADRs dynamically

set -e

# Change to repository root
cd "$(dirname "$0")/.."

echo "🔄 Generating ADR index..."

# Output file
output_file="docs/adr-index.json"

# Start JSON structure
cat > "$output_file" << EOF
{
  "generated": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "version": "1.0", 
  "source": "generate-adr-index.sh",
  "adrs": [
EOF

# Find all ADR files and create array
adr_files=($(find docs/adr -name "[0-9][0-9][0-9][0-9]-*.md" | sort))

# Process each file
for i in "${!adr_files[@]}"; do
  file="${adr_files[$i]}"
  filename=$(basename "$file")
  number=$(echo "$filename" | sed 's/^\([0-9]\{4\}\)-.*/\1/')
  
  echo "  Processing ADR $number: $filename"
  
  # Extract title from first h1 in file (handle quotes properly)
  title=$(grep -m 1 "^# " "$file" | sed 's/^# //' | sed 's/"/\\"/g' || echo "Unknown Title")
  
  # Extract status (look for content after ## Status header)
  status=$(awk '
    /^## Status/ { 
      found_status=1; 
      next 
    } 
    found_status && /^##/ { 
      exit 
    } 
    found_status && NF > 0 && !/^$/ { 
      gsub(/^[ \t]+|[ \t]+$/, "", $0); 
      if($0) { 
        print $0; 
        exit 
      } 
    }
  ' "$file" || echo "Unknown")
  
  # Detect diagram type by looking for C4 diagrams and other types
  diagram_type="-"
  if grep -q "C4Context" "$file"; then
    diagram_type="Context"
  elif grep -q "C4Container" "$file"; then
    diagram_type="Container"
  elif grep -q "C4Component" "$file"; then
    diagram_type="Component"
  elif grep -q -E "(sequenceDiagram|flowchart)" "$file"; then
    diagram_type="Dynamic"
  elif grep -q "\`\`\`mermaid" "$file"; then
    diagram_type="Diagram"
  fi
  
  # Generate relative path for web app (relative to docs/ directory)
  relative_path="adr/$filename"
  
  # Add comma if not first entry
  if [ "$i" -gt 0 ]; then
    echo "," >> "$output_file"
  fi
  
  # Add ADR entry to JSON (properly escape JSON strings)
  cat >> "$output_file" << EOF
    {
      "number": "$number",
      "title": "$title",
      "status": "$status",
      "diagramType": "$diagram_type",
      "filePath": "$relative_path",
      "fileName": "$filename"
    }
EOF
done

# Close JSON array and object
cat >> "$output_file" << 'EOF'
  ]
}
EOF

# Validate JSON and pretty print if possible
if command -v jq >/dev/null 2>&1; then
  echo "📝 Validating and formatting JSON..."
  temp_file=$(mktemp)
  if jq . "$output_file" > "$temp_file"; then
    mv "$temp_file" "$output_file"
    echo "✅ JSON validation successful"
  else
    echo "❌ JSON validation failed"
    rm -f "$temp_file"
    exit 1
  fi
elif command -v python3 >/dev/null 2>&1; then
  echo "📝 Formatting JSON with Python..."
  temp_file=$(mktemp)
  if python3 -m json.tool "$output_file" > "$temp_file"; then
    mv "$temp_file" "$output_file"
    echo "✅ JSON formatting successful"
  else
    echo "❌ JSON formatting failed"
    rm -f "$temp_file"
    exit 1
  fi
else
  echo "📝 Raw JSON generated (no formatter available)"
fi

# Report results
adr_count=$(find docs/adr -name "[0-9][0-9][0-9][0-9]-*.md" | wc -l | tr -d ' ')
echo "✅ Generated ADR index with $adr_count entries"
echo "📄 Output: docs/adr-index.json"

# Show summary if jq is available
if command -v jq >/dev/null 2>&1; then
  echo ""
  echo "📊 Summary:"
  echo "   Total ADRs: $adr_count"
  
  # Count by status
  if [ -f docs/adr-index.json ]; then
    accepted=$(jq -r '.adrs[] | select(.status == "Accepted") | .number' docs/adr-index.json | wc -l | tr -d ' ')
    proposed=$(jq -r '.adrs[] | select(.status == "Proposed") | .number' docs/adr-index.json | wc -l | tr -d ' ')
    deprecated=$(jq -r '.adrs[] | select(.status == "Deprecated") | .number' docs/adr-index.json | wc -l | tr -d ' ')
    superseded=$(jq -r '.adrs[] | select(.status == "Superseded") | .number' docs/adr-index.json | wc -l | tr -d ' ')
    
    echo "   Accepted: $accepted"
    echo "   Proposed: $proposed" 
    echo "   Deprecated: $deprecated"
    echo "   Superseded: $superseded"
    
    # Count diagrams
    with_diagrams=$(jq -r '.adrs[] | select(.diagramType != "-") | .number' docs/adr-index.json | wc -l | tr -d ' ')
    echo "   With Diagrams: $with_diagrams"
  fi
fi