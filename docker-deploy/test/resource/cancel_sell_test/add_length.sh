#!/bin/bash

# Process all XML files in the current directory
for file in *.xml; do
    if [ -f "$file" ]; then
        echo "Processing $file..."
        
        # Calculate the byte length of the XML content
        content=$(cat "$file")
        length=${#content}
        
        # Create a temporary file with the length followed by the content
        echo -e "$length\n$content" > "${file}.tmp"
        
        # Replace the original file
        mv "${file}.tmp" "$file"
        
        echo "Added byte length $length to $file"
    fi
done

echo "All XML files processed successfully."