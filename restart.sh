#!/bin/bash

docker pull your-registry/restful-to-mcp:latest

docker-compose down

docker-compose up -d