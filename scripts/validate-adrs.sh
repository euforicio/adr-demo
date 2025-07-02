#!/bin/bash

# Validate ADR Structure
# This script checks that all ADRs follow the required format

set -e

# Change to repository root
cd "$(dirname "$0")/.."

echo "🔍 Validating ADR structure..."

# Required sections for ADRs
required_sections=("## Status" "## Context" "## Decision" "## Consequences")

# Track validation results
errors=0
warnings=0

# Find all ADR files
find docs/adr -name "[0-9][0-9][0-9][0-9]-*.md" | sort | while IFS= read -r file; do
  filename=$(basename "$file")
  number=$(echo "$filename" | sed 's/^\([0-9]\{4\}\)-.*/\1/')
  
  echo "📄 Validating ADR $number: $filename"
  
  # Check for required sections
  for section in "${required_sections[@]}"; do
    if ! grep -q "^$section" "$file"; then
      echo "  ❌ Missing required section: $section"
      ((errors++))
    fi
  done
  
  # Check for title (h1 header)
  if ! grep -q "^# " "$file"; then
    echo "  ❌ Missing title (h1 header)"
    ((errors++))
  fi
  
  # Check status values
  status=$(awk '/^## Status/ {getline; while(getline && $0 !~ /^##/ && $0 !~ /^$/) {gsub(/^[ \t]+|[ \t]+$/, "", $0); if($0) {print $0; exit}}}' "$file")
  valid_statuses=("Proposed" "Accepted" "Deprecated" "Superseded")
  
  if [[ ! " ${valid_statuses[@]} " =~ " ${status} " ]]; then
    echo "  ⚠️  Invalid status: '$status' (should be one of: ${valid_statuses[*]})"
    ((warnings++))
  fi
  
  # Check for template placeholder text
  if grep -q "Brief description of the decision" "$file"; then
    echo "  ⚠️  Contains template placeholder text"
    ((warnings++))
  fi
  
  # Check filename format
  if [[ ! "$filename" =~ ^[0-9]{4}-[a-z0-9-]+\.md$ ]]; then
    echo "  ❌ Invalid filename format (should be NNNN-kebab-case.md)"
    ((errors++))
  fi
  
  # Check for empty content
  if [ ! -s "$file" ]; then
    echo "  ❌ File is empty"
    ((errors++))
  fi
  
  echo "  ✅ ADR $number validation complete"
done

# Check for sequential numbering
echo ""
echo "🔢 Checking sequential numbering..."

# Get all ADR files sorted by number
adr_files=$(find docs/adr -name '[0-9][0-9][0-9][0-9]-*.md' | sort)

if [ -z "$adr_files" ]; then
  echo "  ✅ No ADR files found"
else
  expected_number=1
  
  for file in $adr_files; do
    filename=$(basename "$file")
    number=$(echo "$filename" | sed 's/^\([0-9][0-9][0-9][0-9]\)-.*/\1/')
    expected_padded=$(printf "%04d" $expected_number)
    
    if [ "$number" != "$expected_padded" ]; then
      echo "  ❌ Expected ADR $expected_padded but found $number"
      echo "     ADRs must be numbered sequentially starting from 0001"
      ((errors++))
    fi
    
    expected_number=$((expected_number + 1))
  done
  
  if [ $expected_number -eq 1 ]; then
    echo "  ✅ No ADR files to validate"
  else
    total_adrs=$((expected_number - 1))
    echo "  ✅ Sequential numbering validated for $total_adrs ADRs"
  fi
fi

# Check for duplicates
echo ""
echo "🔍 Checking for duplicate numbers..."
duplicates=$(printf '%s\n' "${numbers[@]}" | sort | uniq -d)
if [ -n "$duplicates" ]; then
  echo "  ❌ Duplicate ADR numbers found:"
  echo "$duplicates" | sed 's/^/    /'
  ((errors++))
else
  echo "  ✅ No duplicate numbers found"
fi

# Final report
echo ""
echo "📊 Validation Summary"
echo "===================="
echo "Total ADRs: $(find docs/adr -name "[0-9][0-9][0-9][0-9]-*.md" | wc -l | tr -d ' ')"
echo "Errors: $errors"
echo "Warnings: $warnings"

if [ $errors -eq 0 ]; then
  echo "✅ All ADRs pass validation"
  exit 0
else
  echo "❌ Validation failed with $errors errors"
  exit 1
fi