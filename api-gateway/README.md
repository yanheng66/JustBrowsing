# API Gateway

The API Gateway serves as the entry point for all client requests in the Just Browsing e-commerce platform. It implements routing, load balancing, caching, and security features to provide a unified interface for accessing the platform's services.

## Architecture Overview

The API Gateway is built on NGINX and follows the API Gateway pattern in microservices architecture:

### Design Principles

1. **Single Entry Point**: Provides a unified entry point for all client requests
2. **Service Aggregation**: Routes requests to appropriate backend services
3. **Separation of Concerns**: Isolates cross-cutting concerns from service implementations
4. **Caching**: Implements response caching for improved performance
5. **Resilience**: Provides load balancing and failover capabilities
6. **Security**: Centralizes security controls and header management

### Key Features

- Intelligent routing based on URL patterns
- Load balancing across redundant service instances
- Response caching for read operations
- Health check monitoring of backend services
- Security header management
- Detailed access and error logging

## Implementation Details

### Routing Logic

The API Gateway routes requests based on URL patterns:

- **Command Endpoints**: `/api/commands/*` → Command Service
  - Write operations (POST, PUT, DELETE)
  - No caching to ensure data consistency
  
- **Query Endpoints**: `/api/queries/*` → Query Service
  - Read operations (GET)
  - Response caching for improved performance

### Load Balancing

Implements load balancing across redundant service instances:

- **Algorithm**: Least connections
- **Health Checks**: Regular health checks to detect service availability
- **Failover**: Automatic failover to healthy instances

### Caching Strategy

- **Cache Storage**: In-memory cache with configurable size
- **Cache Keys**: Based on request URI
- **TTL**: 5 minutes for successful responses (HTTP 200)
- **Cache Bypass**: Requests with Authorization headers bypass cache
- **Cache Control**: Supports standard HTTP cache control mechanisms

### Security Features

- **HTTP Headers**: Adds security-related headers to responses
  - X-Content-Type-Options
  - X-XSS-Protection
  - X-Frame-Options
  - Content-Security-Policy
  - Strict-Transport-Security
- **TLS Support**: Configurable HTTPS support
- **IP Filtering**: Optional IP allowlist/denylist

## NGINX Configuration

The core of the API Gateway is its NGINX configuration:

```nginx
# Main server configuration
server {
    listen 80;
    listen [::]:80;
    server_name ${SERVER_NAME};

    # Security headers
    add_header X-Content-Type-Options "nosniff";
    add_header X-XSS-Protection "1; mode=block";
    add_header X-Frame-Options "SAMEORIGIN";
    add_header Content-Security-Policy "default-src 'self'";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Health check endpoint
    location = /health {
        access_log off;
        add_header Content-Type application/json;
        return 200 '{"status":"UP","timestamp":"${time_iso8601}"}';
    }

    # Command endpoints
    location /api/commands/ {
        proxy_pass http://command-service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Avoid caching command (write) operations
        proxy_no_cache 1;
        proxy_cache_bypass 1;

        # Timeout settings
        proxy_connect_timeout 5s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Query endpoints
    location /api/queries/ {
        proxy_pass http://query-service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Enable caching for query (read) operations
        proxy_cache api_cache;
        proxy_cache_valid 200 5m;  # Cache successful responses for 5 minutes
        add_header X-Cache-Status $upstream_cache_status;  # Indicate if response was cached
        
        # Skip cache for requests with authorization header
        proxy_cache_bypass $http_authorization;
        
        # Timeout settings
        proxy_connect_timeout 5s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
}
```

## Configuration Options

The API Gateway is designed to be highly configurable through environment variables and configuration files.

### Environment Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `SERVER_NAME` | Server name for NGINX | `localhost` |
| `COMMAND_SERVICE_HOST_1` | Primary Command Service host | `command-service-1` |
| `COMMAND_SERVICE_PORT_1` | Primary Command Service port | `8080` |
| `COMMAND_SERVICE_HOST_2` | Secondary Command Service host | `command-service-2` |
| `COMMAND_SERVICE_PORT_2` | Secondary Command Service port | `8080` |
| `QUERY_SERVICE_HOST_1` | Primary Query Service host | `query-service-1` |
| `QUERY_SERVICE_PORT_1` | Primary Query Service port | `8081` |
| `QUERY_SERVICE_HOST_2` | Secondary Query Service host | `query-service-2` |
| `QUERY_SERVICE_PORT_2` | Secondary Query Service port | `8081` |

### Configuration Files

The API Gateway uses the following configuration files:

- **nginx.conf**: Main NGINX configuration
- **nginx.conf.template**: Template file for environment variable substitution
- **default.env**: Default environment variable values

## Deployment

### Docker

```bash
# Build Docker image
docker build -t justbrowsing/api-gateway .

# Run Docker container
docker run -p 80:80 \
  -e SERVER_NAME=ecommerce.example.com \
  -e COMMAND_SERVICE_HOST_1=command-service-1 \
  -e COMMAND_SERVICE_PORT_1=8080 \
  -e QUERY_SERVICE_HOST_1=query-service-1 \
  -e QUERY_SERVICE_PORT_1=8081 \
  justbrowsing/api-gateway
```

### Docker Compose

```bash
# Start the API Gateway
docker-compose up -d
```

The docker-compose.yml file is configured to:

1. Mount the NGINX configuration files
2. Set environment variables
3. Expose port 80 (and optionally 443 for HTTPS)
4. Configure health checks
5. Set up logging

## Resource Optimization (for t2.micro)

For deployment on t2.micro instances, the following optimizations are configured:

```nginx
# Worker processes and connections
worker_processes auto;  # Automatically set based on available cores
worker_connections 1024;  # Limit connections per worker

# Reduce buffer sizes
client_body_buffer_size 10K;
client_header_buffer_size 1k;
client_max_body_size 8m;
large_client_header_buffers 2 1k;

# Timeouts to prevent resource exhaustion
client_body_timeout 12;
client_header_timeout 12;
keepalive_timeout 15;
send_timeout 10;

# Caching
open_file_cache max=1000 inactive=20s;
open_file_cache_valid 30s;
open_file_cache_min_uses 2;
open_file_cache_errors on;
```

## Logging and Monitoring

### Access Logs

Access logs capture detailed information about requests:

```
log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                '$status $body_bytes_sent "$http_referer" '
                '"$http_user_agent" "$http_x_forwarded_for" '
                '$request_time $upstream_response_time';
```

Access logs are stored in `/var/log/nginx/access.log`.

### Error Logs

Error logs capture NGINX errors and are stored in `/var/log/nginx/error.log`.

### Health Checks

The API Gateway exposes a health check endpoint at `/health` that returns:

```json
{
  "status": "UP",
  "timestamp": "2025-04-23T10:15:30Z"
}
```

This endpoint is used by load balancers to check the gateway's health.

## Operation Guide

### Starting the Gateway

```bash
# Using the start script
./start.sh

# Or directly with Docker Compose
docker-compose up -d
```

### Stopping the Gateway

```bash
docker-compose down
```

### Viewing Logs

```bash
# View access logs
docker-compose logs api-gateway

# Follow logs in real-time
docker-compose logs -f api-gateway

# View specific log files
docker exec -it api-gateway cat /var/log/nginx/access.log
docker exec -it api-gateway cat /var/log/nginx/error.log
```

### Reloading Configuration

```bash
# Reload NGINX configuration without downtime
docker exec -it api-gateway nginx -s reload
```

### Checking Cache Status

The `X-Cache-Status` header in responses indicates cache status:

- `MISS`: Response not cached
- `HIT`: Response served from cache
- `BYPASS`: Cache bypassed (e.g., when Authorization header is present)
- `EXPIRED`: Cache entry expired
- `UPDATING`: Cache entry being updated

## Troubleshooting

### Common Issues

1. **502 Bad Gateway**
   - Check if backend services are running
   - Verify network connectivity to backend services
   - Check service health endpoints

2. **504 Gateway Timeout**
   - Backend service response is too slow
   - Increase timeout values in NGINX configuration

3. **Cache Not Working**
   - Verify cache configuration in NGINX
   - Check if requests include Authorization headers (which bypass cache)
   - Ensure responses have HTTP 200 status code

4. **High CPU/Memory Usage**
   - Reduce worker connections
   - Adjust buffer sizes
   - Check for traffic spikes or DDoS attacks

### Debugging Steps

1. **Check NGINX Configuration**
   ```bash
   docker exec -it api-gateway nginx -t
   ```

2. **Inspect Running Configuration**
   ```bash
   docker exec -it api-gateway nginx -T
   ```

3. **Enable Debug Logging**
   - Set `error_log /var/log/nginx/error.log debug;` in nginx.conf
   - Reload configuration

4. **Check Backend Service Health**
   ```bash
   curl http://localhost/api/commands/actuator/health
   curl http://localhost/api/queries/health
   ```

5. **Verify DNS Resolution**
   ```bash
   docker exec -it api-gateway ping command-service
   docker exec -it api-gateway ping query-service
   ```

## Security Considerations

### TLS Configuration

For production deployment, enable HTTPS with proper TLS configuration:

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;

    ssl_certificate /etc/nginx/ssl/server.crt;
    ssl_certificate_key /etc/nginx/ssl/server.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_dhparam /etc/nginx/ssl/dhparam.pem;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
}
```

### Rate Limiting

Implement rate limiting to prevent abuse:

```nginx
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

location /api/ {
    limit_req zone=api_limit burst=20 nodelay;
    # ...
}
```

### IP Filtering

Restrict access by IP when necessary:

```nginx
# Allow only specific IPs
allow 192.168.1.0/24;
deny all;
```

## Integration with External Services

### CDN Integration

Configure for working with CDNs:

```nginx
# Trust X-Forwarded-For from CDN
set_real_ip_from 192.168.1.0/24;  # CDN IP range
real_ip_header X-Forwarded-For;
real_ip_recursive on;
```

### Authentication Services

Configure for authentication service integration:

```nginx
# Forward requests to authentication service
location /auth/ {
    proxy_pass http://auth-service;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}

# Validate JWT
location /api/ {
    auth_jwt "API";
    auth_jwt_key_file /etc/nginx/jwt_key.pem;
    # ...
}
```