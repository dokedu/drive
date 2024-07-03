#!/bin/bash

# Read the version from the .version file
VERSION=$(cat .version)

# Build and push the Docker image with the read version
docker buildx build --push -t ghcr.io/dokedu/dokedu-drive-backend:$VERSION .
