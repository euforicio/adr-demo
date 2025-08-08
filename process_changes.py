import os
import json
import re

code_changes = os.environ.get('CODE_CHANGES', '')
files_to_create = json.loads(os.environ.get('FILES_TO_CREATE', '[]'))
files_to_modify = json.loads(os.environ.get('FILES_TO_MODIFY', '[]'))

# Split code by file markers
file_pattern = r'--- FILE: (.*?) ---'
parts = re.split(file_pattern, code_changes)

# Process files
if len(parts) > 1:
    # Multiple files specified in the output
    for i in range(1, len(parts), 2):
        if i < len(parts) - 1:
            filepath = parts[i].strip()
            content = parts[i + 1].strip()
            
            # Create directory if needed
            os.makedirs(os.path.dirname(filepath), exist_ok=True) if os.path.dirname(filepath) else None
            
            # Write file
            with open(filepath, 'w') as f:
                f.write(content)
            print(f"Created/Modified: {filepath}")
else:
    # Single file or fallback
    if files_to_create:
        filepath = files_to_create[0]
    elif files_to_modify:
        filepath = files_to_modify[0]
    else:
        filepath = "AI_IMPLEMENTATION.md"
    
    # Create directory if needed
    os.makedirs(os.path.dirname(filepath), exist_ok=True) if os.path.dirname(filepath) else None
    
    # Write the content
    with open(filepath, 'w') as f:
        f.write(code_changes)
    print(f"Created/Modified: {filepath}")
