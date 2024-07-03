#!/bin/bash

# Function to increment the version number
increment_version() {
    local version=$1
    local base_version=${version%.*}
    local patch_version=${version##*.}
    local new_patch_version=$((patch_version + 1))
    echo "${base_version}.${new_patch_version}"
}

# Increment the version in the specified directory
increment_version_in_dir() {
    local dir=$1
    local version_file="${dir}/.version"

    if [ -f "$version_file" ]; then
        local current_version=$(cat "$version_file")
        local new_version=$(increment_version "$current_version")
        echo "$new_version" > "$version_file"
        echo "Updated $dir to version $new_version"
    else
        echo "Version file not found in $dir"
    fi
}

# Increment versions in both frontend and backend directories
increment_version_in_dir "frontend"
increment_version_in_dir "backend"
