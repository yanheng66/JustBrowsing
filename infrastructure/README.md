# JustBrowsing Infrastructure

This directory contains all the infrastructure configuration for the JustBrowsing e-commerce platform. The infrastructure is designed to support the CQRS architecture with separate Command and Query services.

## Directory Structure

- **docker/** - Docker Compose files for local development
- **monitoring/** - Monitoring configuration (Prometheus, Grafana, ELK, Jaeger)
- **terraform/** - Terraform files for AWS deployment

## Local Development Environment

The local development environment is set up using Docker Compose and includes all necessary services:

### Core Services
- Command Service (Spring Boot)
- Query Service (Go)
- API Gateway (NGINX)

### Databases
- PostgreSQL (for Command Service)
- MongoDB (for Query Service)
- Redis (for caching)

### Messaging
- Kafka and Zookeeper (for event streaming)
- Schema Registry (for Avro schemas)

### Search
- Elasticsearch (for product search)

### Monitoring
- Prometheus (for metrics collection)
- Grafana (for dashboards)
- ELK Stack (for log aggregation)
- Jaeger (for distributed tracing)

## Getting Started

### Starting the Local Development Environment

To start the entire local development environment, run:

```bash
cd infrastructure/docker
./start-local-dev.sh
```

This script will:
1. Create necessary configuration files
2. Start all required services
3. Initialize databases with required schemas
4. Set up monitoring tools

### Accessing Services

- Application: http://localhost
- Grafana: http://localhost:3000 (admin/admin)
- Kibana: http://localhost:5601
- Jaeger UI: http://localhost:16686
- Prometheus: http://localhost:9090

## AWS Deployment

For AWS deployment, Terraform is used to provision the required infrastructure. The deployment follows the architecture specified in the design document.

### Deployment Architecture

All instances use t2.micro due to AWS account limitations. The services are consolidated to maximize resource utilization:

| EC2 Instance | Services                                 | Subnet Type |
| ------------ | ---------------------------------------- | ----------- |
| Instance 1   | API Gateway (NGINX), Prometheus, Grafana | Public      |
| Instance 2   | API Gateway (NGINX) - Redundancy         | Public      |
| Instance 3   | Command Service, Kafka                   | Private     |
| Instance 4   | Command Service - Redundancy             | Private     |
| Instance 5   | Query Service, Redis                     | Private     |
| Instance 6   | Query Service - Redundancy               | Private     |
| Instance 7   | PostgreSQL, ELK Stack                    | Private     |
| Instance 8   | MongoDB, Elasticsearch                   | Private     |

### Load Balancing

AWS Elastic Load Balancer (ELB) configuration:

- **Public-facing ELB**:
  - Forwards traffic to the API Gateway instances (Instances 1 and 2)
  - HTTP on port 80, HTTPS on port 443
  - Health check path: `/health`
  - Health check interval: 30 seconds
- **Internal Command ELB**:
  - Routes traffic from API Gateway to Command Service instances (Instances 3 and 4)
  - Health check path: `/actuator/health`
- **Internal Query ELB**:
  - Routes traffic from API Gateway to Query Service instances (Instances 5 and 6)
  - Health check path: `/health`

### Resource Optimization

Given the t2.micro constraints, the following optimizations are configured:

1. **JVM Settings for Spring Boot**:
   - Maximum heap size: 512MB
   - Compressed object pointers
   - Reduced thread stack size
2. **Go Service Settings**:
   - GOMAXPROCS set to match available vCPUs (1 for t2.micro)
   - Limit concurrent connections to prevent resource exhaustion
3. **Database Configurations**:
   - PostgreSQL: max_connections=50, shared_buffers=128MB
   - MongoDB: limit WiredTiger cache to 256MB
   - Elasticsearch: heap size of 512MB, disable ML features
4. **Docker Container Configurations**:
   - Memory limits enforced for all containers
   - CPU share limits to prevent single container monopolizing resources
   - Use of Alpine-based images to minimize footprint

### Deploying to AWS

To deploy to AWS, follow these steps:

1. **Prerequisites**:
   - AWS CLI configured with appropriate credentials
   - An SSH key pair created in your AWS account
   - An S3 bucket and DynamoDB table for Terraform state management

2. **Initialize Terraform**:
   ```bash
   cd infrastructure/terraform/main
   terraform init
   ```

3. **Plan the deployment**:
   ```bash
   terraform plan -var="key_name=your-key-name"
   ```

4. **Apply the configuration**:
   ```bash
   terraform apply -var="key_name=your-key-name"
   ```

5. **To destroy the infrastructure**:
   ```bash
   terraform destroy -var="key_name=your-key-name"
   ```

## Monitoring and Observability

The monitoring infrastructure includes:

1. **Health Checks**:
   - Command Service: `/actuator/health` endpoint exposing service health
   - Query Service: `/health` endpoint providing service status
   - API Gateway: `/health` endpoint for ELB health checks

2. **Metrics Collection**:
   - Prometheus scrapes metrics from all services
   - Grafana dashboards for visualization
   - JVM metrics for Command Service
   - Go runtime metrics for Query Service
   - Database connection and performance metrics

3. **Log Aggregation**:
   - ELK Stack (Elasticsearch, Logstash, Kibana)
   - Filebeat for log collection
   - Centralized log storage and search

4. **Distributed Tracing**:
   - Jaeger for tracing request flows
   - Visualization of cross-service communication