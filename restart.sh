#!/bin/bash

docker push your-registry/restful-to-mcp:latest

docker-compose down
docker-compose up -d