#!/bin/bash

# Function to prompt for deployment confirmation
prompt_deploy() {
    read -p "Do you want to deploy? (y/n): " confirm
    if [[ $confirm == "y" || $confirm == "Y" ]]; then
        return 0
    else
        return 1
    fi
}

# Main script execution
if prompt_deploy; then
    echo "Deploying..."

    # Call the script to increment the version
    ./scripts/version.sh
    if [ $? -ne 0 ]; then
        echo "Version increment failed!"
        exit 1
    fi

    # Call the script to build both frontend and backend
    ./scripts/build.sh
    if [ $? -ne 0 ]; then
        echo "Build failed!"
        exit 1
    fi

    echo "Deployment successful!"
else
    echo "Deployment cancelled."
fi
