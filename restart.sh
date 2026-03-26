#!/bin/bash

# Replace with your own image registry address
IMAGE="your-registry/restful-to-mcp:latest"

docker pull $IMAGE

docker-compose down

docker-compose up -d