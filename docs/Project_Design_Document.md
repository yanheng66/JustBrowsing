# Just Browsing Design Document

## 1. System Overview

### 1.1 Project Background

This document outlines the technical design of a lightweight e-commerce platform based on the CQRS (Command Query Responsibility Segregation) architectural pattern. The system employs an event-driven architecture to separate write operations (command side) and read operations (query side). The design features asynchronous updates on the query side (eventual consistency), while the write side processes critical operations (such as order placement and inventory management) synchronously to enforce strict inventory validation.

### 1.2 Core Problem Solution

The main challenge for an e-commerce platform is ensuring data consistency while providing high performance and scalability. Particularly in inventory management, the system must prevent overselling while handling a large volume of concurrent query requests. By implementing CQRS, we can optimize each side of the system independently:

- **Command Side**: Optimized for data consistency and business rule enforcement
- **Query Side**: Optimized for read performance and search capabilities

## 2. Technology Stack

### 2.1 Command Service (Write Side)

- **Framework**: Spring Boot
- **Database**: PostgreSQL
- **Data Access**: Spring Data JPA/Hibernate
- **Transaction Management**: Spring Transaction Management (local transactions)

### 2.2 Query Service (Read Side)

- **Framework**: Go with Gin
- Data Storage:
  - MongoDB (primary data storage)
  - Elasticsearch (search functionality)
  - Redis (caching)
- **Query Optimization**: Strategic indexing and data denormalization

### 2.3 Event Bus/Messaging System

- **Message Queue**: Apache Kafka
- **Serialization Format**: Avro/Protocol Buffers
- **Event Storage**: Kafka as event log

### 2.4 API Gateway

- NGINX

### 2.5 Deployment and Operations

- **Containerization**: Docker
- **Monitoring**: Prometheus + Grafana + ELK Stack + Jaeger
- **Infrastructure as Code**: Packer + Terraform

## 3. Deployment Architecture

### 3.1 Instance Allocation

All instances will be t2.micro due to AWS account limitations. The services will be consolidated to maximize resource utilization:

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

### 3.2 Load Balancing

AWS Elastic Load Balancer (ELB) configuration:

- Public-facing ELB:
  - Forwards traffic to the API Gateway instances (Instances 1 and 2)
  - HTTP on port 80, HTTPS on port 443
  - Health check path: `/health`
  - Health check interval: 30 seconds
- Internal Command ELB:
  - Routes traffic from API Gateway to Command Service instances (Instances 3 and 4)
  - Health check path: `/actuator/health`
- Internal Query ELB:
  - Routes traffic from API Gateway to Query Service instances (Instances 5 and 6)
  - Health check path: `/health`

### 3.3 Resource Optimization

Given the t2.micro constraints, the following optimizations are necessary:

1. JVM Settings for Spring Boot:
   - Maximum heap size: 512MB
   - Compressed object pointers
   - Reduced thread stack size
2. Go Service Settings:
   - GOMAXPROCS set to match available vCPUs (1 for t2.micro)
   - Limit concurrent connections to prevent resource exhaustion
3. Database Configurations:
   - PostgreSQL: max_connections=50, shared_buffers=128MB
   - MongoDB: limit WiredTiger cache to 256MB
   - Elasticsearch: heap size of 512MB, disable ML features
4. Docker Container Configurations:
   - Memory limits enforced for all containers
   - CPU share limits to prevent single container monopolizing resources
   - Use of Alpine-based images to minimize footprint

## 4. Detailed Architecture Design

The architecture follows a strict CQRS pattern with event sourcing to propagate changes from the command side to the query side.

### 4.1 Command Service (Write Side)

#### 4.1.1 Component Structure

- **Controller Layer**: REST API endpoints for commands
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data access abstractions
- **Domain Model**: Core business entities
- **Event Publisher**: Produces events to Kafka

#### 4.1.2 Process Flow

1. Command request received via API
2. Input validation
3. Business logic execution within transaction boundary
4. Database updates
5. Event publishing on successful transaction
6. Response to client

#### 4.1.3 Database Design (PostgreSQL)

**Products Table**

```sql
products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
```

**Tags Table**

```sql
tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
```

**ProductTags Table**

```sql
product_tags (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    tag_id INTEGER REFERENCES tags(id),
    tag_value VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (product_id, tag_id)
)
```

**Inventory Table**

```sql
inventory (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 0,
    last_replenishment_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_quantity CHECK (quantity >= 0)
)
```

**Orders Table**

```sql
orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
```

**OrderItems Table**

```sql
order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    CONSTRAINT positive_quantity CHECK (quantity > 0)
)
```

**Outbox Table (for transactional outbox pattern)**

```sql
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

#### 4.1.4 Optimization Strategies

- **Optimistic Locking**: Using version fields for concurrent inventory updates
- **Transactional Outbox Pattern**: Ensures reliable event publishing
- Database Indexes:
  - B-tree index on products(sku)
  - B-tree index on inventory(product_id)
  - B-tree index on orders(order_number)
  - B-tree index on order_items(order_id)
  - B-tree index on outbox_events(processed, created_at) for efficient event processing
  - B-tree index on product_tags(product_id, tag_id)
  - B-tree index on tags(name)

### 4.2 Query Service (Read Side)

#### 4.2.1 Component Structure

- **API Routes**: REST endpoints for queries
- **Service Layer**: Query execution and data retrieval
- **Repository Layer**: Data access abstractions
- **Event Consumers**: Process events from Kafka
- **Cache Manager**: Handle Redis caching

#### 4.2.2 Data Models

**MongoDB Collections**

**Products Collection**

```json
{
  "_id": ObjectId,
  "productId": String,
  "sku": String,
  "name": String,
  "description": String,
  "price": Decimal128,
  "tags": [
    {
      "name": String,
      "value": String
    }
  ],
  "currentInventory": Number,
  "images": [String],
  "created": Date,
  "updated": Date
}
```

**Orders Collection**

```json
{
  "_id": ObjectId,
  "orderId": String,
  "orderNumber": String,
  "totalAmount": Decimal128,
  "items": [{
    "productId": String,
    "productName": String,
    "sku": String,
    "quantity": Number,
    "unitPrice": Decimal128,
    "totalPrice": Decimal128
  }],
  "created": Date
}
```

**Elasticsearch Mappings**

**Products Index**

```json
{
  "mappings": {
    "properties": {
      "productId": { "type": "keyword" },
      "sku": { "type": "keyword" },
      "name": { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
      "description": { "type": "text" },
      "price": { "type": "float" },
      "tags": {
        "properties": {
          "name": { "type": "keyword" },
          "value": { "type": "keyword" }
        }
      },
      "currentInventory": { "type": "integer" },
      "created": { "type": "date" },
      "updated": { "type": "date" }
    }
  }
}
```

**Redis Cache Keys**

- `product:{productId}` - Product details
- `products:tag:{tagName}:{tagValue}` - List of products by tag
- `inventory:{productId}` - Current inventory level

#### 4.2.3 Event Handling

The Query Service consumes events from Kafka and updates the read models accordingly:

- ProductCreated → Add to MongoDB, Elasticsearch, invalidate caches
- ProductUpdated → Update MongoDB, Elasticsearch, invalidate caches
- ProductTagAdded → Update product tags in MongoDB, update Elasticsearch, invalidate caches
- ProductTagRemoved → Remove product tag from MongoDB and Elasticsearch, invalidate caches
- InventoryChanged → Update inventory in MongoDB, update Redis cache
- OrderCreated → Add to MongoDB

#### 4.2.4 Optimization Strategies

- MongoDB Indexes:
  - productId, sku (unique)
  - orderId, orderNumber (unique)
  - "tags.name", "tags.value" (for tag searches)
- Elasticsearch Optimizations:
  - Custom analyzers for product search
  - Keyword fields for exact matches and aggregations
  - Optimized indexes on tag fields
- Redis Caching Strategy:
  - Time-based expiration for product details (TTL: 1 hour)
  - Invalidation on relevant events
  - Write-through caching for high-frequency queries

## 5. Key Interface Design

### 5.1 Command API Endpoints

#### 5.1.1 Product Management

- Create Product
  - `POST /api/commands/products`
  - Request: Product details (name, price, etc.)
  - Response: Product ID, status
- Update Product
  - `PUT /api/commands/products/{productId}`
  - Request: Updated product details
  - Response: Status
- Add Product Tag
  - `POST /api/commands/products/{productId}/tags`
  - Request: Tag name and value
  - Response: Status
- Remove Product Tag
  - `DELETE /api/commands/products/{productId}/tags/{tagId}`
  - Response: Status
- Update Inventory
  - `PUT /api/commands/products/{productId}/inventory`
  - Request: Quantity change (positive for restock, negative for adjustments)
  - Response: New inventory level, status

#### 5.1.2 Order Management

- Create Order
  - `POST /api/commands/orders`
  - Request: Order items
  - Response: Order ID, order number, status

### 5.2 Query API Endpoints

#### 5.2.1 Product Queries

- Get Product
  - `GET /api/queries/products/{productId}`
  - Response: Complete product details
- Search Products by Tags
  - `GET /api/queries/products/search?tags={tagName1}:{tagValue1},{tagName2}:{tagValue2}`
  - Response: Product list matching all specified tags

### 5.3 API Gateway Configuration

NGINX configuration for routing and load balancing:

```nginx
# Command endpoints
location /api/commands/ {
    proxy_pass http://command-service;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}

# Query endpoints
location /api/queries/ {
    proxy_pass http://query-service;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_cache api_cache;
    proxy_cache_valid 200 5m;  # Cache successful responses for 5 minutes
}
```

## 6. Data Flow and Event Processing

### 6.1 Order Creation Flow

1. Client sends order creation request
2. Command service begins transaction:
   - Validates inventory availability for each item
   - Creates order record
   - Creates order items
   - Decrements inventory for each product
   - Stores OrderCreated event in outbox table
   - Commits transaction
3. Outbox processor publishes events to Kafka
4. Query service consumes the event:
   - Updates MongoDB with new order information
   - Updates product inventory levels in MongoDB and Elasticsearch
   - Invalidates relevant Redis caches

### 6.2 Product Search Flow

1. Client sends tag-based search request
2. Query service processes request:
   - Parses tag parameters from the request
   - Checks Redis cache for query results
   - If cache miss, queries Elasticsearch
   - Filters products containing all requested tags
   - Formats and returns results
   - Updates cache for future requests

### 6.3 Inventory Update Flow

1. Admin/system sends inventory update request
2. Command service begins transaction:
   - Retrieves product inventory with optimistic locking
   - Updates inventory quantity
   - Stores InventoryChanged event in outbox table
   - Commits transaction
3. Outbox processor publishes event to Kafka
4. Query service consumes the event:
   - Updates product inventory in MongoDB
   - Updates inventory in Elasticsearch
   - Invalidates relevant Redis caches

## 7. Performance Considerations

### 7.1 Throughput Optimization

- **Connection Pooling**: Configured for databases to efficiently reuse connections
- **Batch Processing**: Events are processed in batches to improve throughput
- **Asynchronous Processing**: Non-critical operations are handled asynchronously

### 7.2 Latency Optimization

- **Caching Strategy**: Frequent queries cached in Redis
- **Database Indexing**: Strategic indexes to improve query performance
- **Read Models**: Denormalized data structures for efficient queries

### 7.3 Resilience Patterns

- **Circuit Breaker**: Preventing cascading failures between services
- **Retry with Backoff**: Automatic retries for transient failures

## 8. Monitoring and Observability

### 8.1 Health Checks

- **Command Service**: `/actuator/health` endpoint exposing service health
- **Query Service**: `/health` endpoint providing service status
- **Databases**: Custom probes to verify connection and query execution

### 8.2 Metrics Collection

- **Service Metrics**: Request counts, response times, error rates
- **JVM Metrics**: Memory usage, garbage collection, thread counts
- **Go Metrics**: Goroutine counts, memory allocation, GC pauses
- **Database Metrics**: Connection pool stats, query performance