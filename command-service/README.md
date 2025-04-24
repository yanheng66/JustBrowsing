# Command Service

The Command Service is the write-side component of the Just Browsing e-commerce platform, implementing the Command Query Responsibility Segregation (CQRS) pattern. This service is responsible for handling all operations that modify the system state, including product management, inventory tracking, and order processing.

## Architecture Overview

The Command Service is built on Spring Boot and follows a layered architecture:

1. **Controller Layer**: REST API endpoints for commands
2. **Service Layer**: Business logic implementation
3. **Repository Layer**: Data access abstractions using Spring Data JPA
4. **Domain Model**: Core business entities

### Key Features

- Transactional processing of write operations
- Optimistic locking for concurrent inventory updates
- Event publishing via the Transactional Outbox pattern
- Robust validation and error handling

## Domain Models

The service defines the following core domain entities:

- **Product**: Represents a product in the system with attributes like SKU, name, description, price, etc.
- **Tag**: Represents product categorization and attributes
- **ProductTag**: Maps tags to products with specific values
- **Inventory**: Tracks product stock levels with optimistic locking for concurrency
- **Order and OrderItem**: Represent customer orders and the products they contain
- **OutboxEvent**: Implements the Transactional Outbox pattern for reliable event publishing

## Interaction with Other Services

### Event Publishing

The Command Service publishes events to Kafka using the Transactional Outbox pattern:

1. Write operations are performed in a transaction
2. Events are stored in the outbox table as part of the same transaction
3. A scheduled task (OutboxProcessorService) polls for unprocessed events and publishes them to Kafka
4. The Query Service consumes these events to update its read models

### Event Types

- **ProductCreated**: When a new product is added
- **ProductUpdated**: When a product's information is modified
- **ProductTagAdded**: When a tag is added to a product
- **ProductTagRemoved**: When a tag is removed from a product
- **InventoryUpdated**: When a product's inventory is modified
- **OrderCreated**: When a new order is placed

## API Endpoints

### Product Management

- **Create Product**: `POST /api/commands/products`
- **Update Product**: `PUT /api/commands/products/{productId}`
- **Add Product Tag**: `POST /api/commands/products/{productId}/tags`
- **Remove Product Tag**: `DELETE /api/commands/products/{productId}/tags/{tagId}`
- **Update Inventory**: `PUT /api/commands/products/{productId}/inventory`

### Order Management

- **Create Order**: `POST /api/commands/orders`

## Database Schema

The service uses PostgreSQL with the following schema:

```
products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)

tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)

product_tags (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    tag_id INTEGER REFERENCES tags(id),
    tag_value VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (product_id, tag_id)
)

inventory (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 0,
    last_replenishment_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_quantity CHECK (quantity >= 0)
)

orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)

order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    CONSTRAINT positive_quantity CHECK (quantity > 0)
)

outbox_events (
    id SERIAL PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMP WITH TIME ZONE
)
```

## Configuration

### application.properties

```properties
# Application
spring.application.name=command-service
server.port=8080
server.servlet.context-path=/api/commands

# Database Configuration
spring.datasource.url=jdbc:postgresql://localhost:5432/ecommerce
spring.datasource.username=postgres
spring.datasource.password=postgres
spring.datasource.driver-class-name=org.postgresql.Driver

# JPA Configuration
spring.jpa.hibernate.ddl-auto=update
spring.jpa.properties.hibernate.dialect=org.hibernate.dialect.PostgreSQLDialect
spring.jpa.show-sql=true
spring.jpa.properties.hibernate.format_sql=true

# Kafka Configuration
spring.kafka.bootstrap-servers=localhost:9092
spring.kafka.producer.key-serializer=org.apache.kafka.common.serialization.StringSerializer
spring.kafka.producer.value-serializer=io.confluent.kafka.serializers.KafkaAvroSerializer
spring.kafka.producer.properties.schema.registry.url=http://localhost:8081

# Outbox Processor Configuration
outbox.polling.interval.ms=1000
outbox.max-items-per-polling=100

# Actuator Configuration
management.endpoints.web.exposure.include=health,info,metrics
management.endpoint.health.show-details=always

# Logging Configuration
logging.level.org.springframework=INFO
logging.level.com.ecommerce.command=DEBUG
logging.level.org.hibernate.SQL=DEBUG
logging.level.org.hibernate.type.descriptor.sql.BasicBinder=TRACE
```

### Resource Optimization (for t2.micro)

For deployment on t2.micro instances, the following JVM options are recommended:

```
JAVA_OPTS="-Xms256m -Xmx512m -XX:+UseCompressedOops -Xss256k"
```

## Build and Run

### Prerequisites

- Java 17
- Maven
- PostgreSQL
- Kafka and Schema Registry

### Build

```bash
# Build the application
mvn clean package

# Build skipping tests
mvn clean package -DskipTests
```

### Run

```bash
# Run with Maven
mvn spring-boot:run

# Run as JAR file
java -jar target/command-service-0.0.1-SNAPSHOT.jar

# Run with custom configuration
java -jar target/command-service-0.0.1-SNAPSHOT.jar \
  --spring.datasource.url=jdbc:postgresql://custom-host:5432/ecommerce \
  --spring.kafka.bootstrap-servers=custom-kafka:9092
```

### Docker

```bash
# Build Docker image
docker build -t justbrowsing/command-service .

# Run Docker container
docker run -p 8080:8080 \
  -e SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/ecommerce \
  -e SPRING_DATASOURCE_USERNAME=postgres \
  -e SPRING_DATASOURCE_PASSWORD=postgres \
  -e SPRING_KAFKA_BOOTSTRAP_SERVERS=kafka:9092 \
  -e SPRING_KAFKA_PRODUCER_PROPERTIES_SCHEMA_REGISTRY_URL=http://schema-registry:8081 \
  -e JAVA_OPTS="-Xms256m -Xmx512m -XX:+UseCompressedOops -Xss256k" \
  justbrowsing/command-service
```

## Monitoring

The service exposes metrics and health information via Spring Boot Actuator:

- Health endpoint: `/api/commands/actuator/health`
- Metrics endpoint: `/api/commands/actuator/metrics`
- Prometheus endpoint: `/api/commands/actuator/prometheus`

## Development Notes

### Adding a New Event Type

1. Create a new Avro schema in `src/main/resources/avro/`
2. Run `mvn generate-sources` to generate Java classes from the schema
3. Update the relevant service to publish events using the new schema
4. Ensure the Query Service has a corresponding handler

### Handling Concurrency

The service uses optimistic locking for inventory updates:

1. The `Inventory` entity has a `version` field annotated with `@Version`
2. When concurrent updates occur, JPA will throw an `OptimisticLockException`
3. The service catches and handles this exception, typically by retrying the operation

### Transaction Boundaries

All service methods that modify data are annotated with `@Transactional` to ensure ACID properties. This is especially important for the Outbox pattern, where data changes and event creation must occur atomically.