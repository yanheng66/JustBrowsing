#!/bin/bash

# Script to start the API Gateway

# Ensure network exists
echo "Creating Docker network if it doesn't exist..."
docker network create ecommerce-network 2>/dev/null || true

# Start the API Gateway
echo "Starting API Gateway..."
docker-compose up -d

# Check if gateway is running
echo "Checking API Gateway status..."
sleep 5
if docker ps | grep -q api-gateway; then
  echo "API Gateway is running."
  echo "Health check endpoint: http://localhost/health"
  echo "Command service endpoint: http://localhost/api/commands/"
  echo "Query service endpoint: http://localhost/api/queries/"
else
  echo "Error: API Gateway failed to start."
  docker-compose logs api-gateway
fi