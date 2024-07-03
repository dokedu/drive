#!/bin/bash

# Navigate to the backend directory and run the build script
echo "Building backend..."
cd backend
./scripts/build.sh
if [ $? -ne 0 ]; then
    echo "Backend build failed!"
    exit 1
fi
cd ..

# Navigate to the frontend directory and run the build script
echo "Building frontend..."
cd frontend
./scripts/build.sh
if [ $? -ne 0 ]; then
    echo "Frontend build failed!"
    exit 1
fi
cd ..

echo "Both frontend and backend built successfully!"
