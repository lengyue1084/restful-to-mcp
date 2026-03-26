#!/bin/bash

# Replace with your own image registry address
IMAGE="your-registry/restful-to-mcp:latest"

# Build
docker build -t $IMAGE .

# Push (login first: docker login your-registry)
docker push $IMAGE