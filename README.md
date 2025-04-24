# Just Browsing - CQRS E-Commerce Platform

Just Browsing is a lightweight e-commerce platform based on the Command Query Responsibility Segregation (CQRS) architectural pattern. The system employs an event-driven architecture to separate write operations (command side) and read operations (query side), providing optimized performance and scalability.

## System Overview

The platform consists of the following core components:

1. **Command Service** (Java/Spring Boot): Handles all write operations and business logic
2. **Query Service** (Go/Gin): Optimized for read operations and search functionality
3. **API Gateway** (NGINX): Routes requests and provides caching and security features
4. **Event Bus** (Kafka): Enables asynchronous communication between services
5. **Databases**:
   - PostgreSQL: Transactional database for the Command Service
   - MongoDB: Document database for the Query Service's read models
   - Elasticsearch: Search engine for product queries
   - Redis: Caching layer for frequently accessed data

## Architecture Highlights

- **CQRS Pattern**: Separate write and read models for optimized performance
- **Event-Driven Design**: Loosely coupled services communicating via events
- **Eventual Consistency**: Asynchronous updates to the read models
- **Optimistic Concurrency**: Prevents conflicts in inventory management
- **Transactional Outbox**: Ensures reliable event publishing
- **Multi-tier Caching**: Redis and NGINX caching for improved performance
- **Resilience Patterns**: Circuit breakers, retries, and fallbacks

## Key Features

- Product management (create, update, tag, inventory)
- Order processing with inventory validation
- Tag-based product search
- Response caching for improved performance
- Distributed tracing and monitoring
- Horizontal scaling with load balancing
- Resource-optimized for cloud deployment

## Directory Structure

```
JustBrowsing/
├── api-gateway/               # NGINX-based API Gateway
├── command-service/           # Spring Boot Command Service
├── docs/                      # Design and API documentation
├── infrastructure/            # Infrastructure configuration
│   ├── docker/                # Docker Compose for local development
│   ├── monitoring/            # Prometheus, Grafana, ELK, Jaeger
│   └── terraform/             # AWS deployment with Terraform
└── query-service/             # Go-based Query Service
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Java 17
- Go 1.19+
- Maven
- AWS CLI (for cloud deployment)
- Terraform (for cloud deployment)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/JustBrowsing.git
   cd JustBrowsing
   ```

2. Start the local development environment:
   ```bash
   cd infrastructure/docker
   ./start-local-dev.sh
   ```

   This will start all necessary services:
   - PostgreSQL, MongoDB, Redis, Elasticsearch
   - Kafka and Schema Registry
   - Command Service and Query Service
   - API Gateway
   - Monitoring tools (Prometheus, Grafana, ELK, Jaeger)

3. Access the application:
   - Application: http://localhost
   - Grafana: http://localhost:3000 (admin/admin)
   - Kibana: http://localhost:5601
   - Jaeger UI: http://localhost:16686
   - Prometheus: http://localhost:9090

### Building Individual Services

#### Command Service
```bash
cd command-service
mvn clean package
```

#### Query Service
```bash
cd query-service
go build -o query-service
```

## API Endpoints

### Command API

- **Create Product**: `POST /api/commands/products`
- **Update Product**: `PUT /api/commands/products/{productId}`
- **Add Product Tag**: `POST /api/commands/products/{productId}/tags`
- **Remove Product Tag**: `DELETE /api/commands/products/{productId}/tags/{tagId}`
- **Update Inventory**: `PUT /api/commands/products/{productId}/inventory`
- **Create Order**: `POST /api/commands/orders`

### Query API

- **Get Product**: `GET /api/queries/products/{productId}`
- **Search Products by Tags**: `GET /api/queries/products/search?tags={tagName1}:{tagValue1},{tagName2}:{tagValue2}`
- **Get Order**: `GET /api/queries/orders/{orderId}`

## Deployment

### AWS Deployment

The platform is designed to be deployed to AWS using Terraform. The deployment follows the architecture specified in the design document, using t2.micro instances with optimized resource allocation:

```bash
cd infrastructure/terraform/main
terraform init
terraform apply -var="key_name=your-key-name"
```

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

## Monitoring and Observability

The platform includes comprehensive monitoring and observability features:

1. **Metrics**: Prometheus collects metrics from all services
2. **Dashboards**: Grafana provides visualization of metrics
3. **Logs**: ELK Stack centralizes logs from all components
4. **Tracing**: Jaeger enables distributed tracing across services
5. **Health Checks**: All services expose health endpoints

## Development Guidelines

### Adding a New Feature

1. Identify whether the feature is a command (write) or query (read) operation
2. Update the appropriate service (Command or Query)
3. Add new events if needed for cross-service communication
4. Update the API Gateway configuration if necessary
5. Add monitoring and appropriate test coverage

### Testing

Each service has its own testing approach:

- **Command Service**: JUnit tests for Java components
- **Query Service**: Go tests for business logic
- **Integration Tests**: Tests across service boundaries
- **Load Tests**: Performance testing under load

## Documentation

For more detailed information, refer to:

- [Project Design Document](docs/Project_Design_Document.md): Detailed architecture and design decisions
- [API Documentation](docs/API_Documentation.md): Comprehensive API specification
- [Command Service README](command-service/README.md): Command Service details
- [Query Service README](query-service/README.md): Query Service details
- [API Gateway README](api-gateway/README.md): API Gateway configuration
- [Infrastructure README](infrastructure/README.md): Infrastructure and deployment

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- This project implements the CQRS pattern as described by Martin Fowler and Greg Young
- Inspiration for the event-driven architecture from various microservices resources
- Optimized for educational purposes to demonstrate modern distributed system patterns