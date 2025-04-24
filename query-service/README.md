# Query Service

The Query Service is the read-side component of the Just Browsing e-commerce platform, implementing the Command Query Responsibility Segregation (CQRS) pattern. This service is responsible for handling all read operations, providing optimized data models for querying product information and orders.

## Architecture Overview

The Query Service is built with Go and the Gin framework, following a clean architecture approach:

1. **API Layer**: REST endpoints for queries using Gin
2. **Service Layer**: Business logic for data retrieval and query execution
3. **Repository Layer**: Data access abstractions for MongoDB
4. **Event Consumers**: Processing events from Kafka to update read models
5. **Cache Manager**: Handling Redis caching for performance optimization

### Key Features

- Denormalized data models for efficient querying
- Multi-tier caching strategy with Redis
- Full-text search capabilities using Elasticsearch
- Event-driven updates from the Command Service via Kafka
- High-performance, concurrent request handling

## Data Models

The service uses the following data models stored in MongoDB:

### Products Collection

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

### Orders Collection

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

## Interaction with Other Services

### Event Consumption

The Query Service consumes events from Kafka to keep its read models up-to-date:

1. The EventService subscribes to product and order topics
2. When events are received, they are deserialized and processed
3. The appropriate service method is called to update the data models
4. Caches are invalidated or updated as needed

### Event Types Handled

- **ProductCreated**: Creates a new product in MongoDB and Elasticsearch
- **ProductUpdated**: Updates an existing product's information
- **ProductTagAdded**: Adds a tag to a product
- **ProductTagRemoved**: Removes a tag from a product
- **InventoryUpdated**: Updates a product's inventory level
- **OrderCreated**: Creates a new order record

## API Endpoints

### Product Queries

- **Get Product**: `GET /api/queries/products/{productId}`
- **Search Products by Tags**: `GET /api/queries/products/search?tags={tagName1}:{tagValue1},{tagName2}:{tagValue2}`

### Order Queries

- **Get Order**: `GET /api/queries/orders/{orderId}`

## Caching Strategy

The service implements a multi-tier caching strategy:

1. **Product Details Caching**
   - Key: `product:{productId}`
   - TTL: 1 hour
   - Invalidated on product updates

2. **Search Results Caching**
   - Key: `products:tags:{tags-query-string}`
   - TTL: 1 hour
   - Invalidated on product or tag changes

3. **Inventory Level Caching**
   - Key: `inventory:{productId}`
   - TTL: 5 minutes
   - Updated on inventory changes

## Search Functionality

Product search is implemented using Elasticsearch with:

- Custom analyzer for product search
- Keyword fields for exact matches
- Nested queries for tag filtering

If Elasticsearch is unavailable, the service falls back to MongoDB queries.

## Configuration

The service uses a YAML configuration file (config.yaml):

```yaml
server:
  port: 8081
  basePath: "/api/queries"

mongodb:
  uri: "mongodb://localhost:27017"
  database: "ecommerce"
  poolSize: 100
  timeout: 30

elasticsearch:
  addresses: ["http://localhost:9200"]
  username: ""
  password: ""
  indexPrefix: "ecommerce_"

redis:
  address: "localhost:6379"
  password: ""
  db: 0
  poolSize: 10
  ttl: 3600  # 1 hour in seconds

kafka:
  brokers: ["localhost:9092"]
  groupId: "query-service"
  topics:
    product: "products"
    inventory: "inventory"
    order: "orders"

logging:
  level: "debug"  # debug, info, warn, error
  format: "json"  # json or text
```

### Resource Optimization (for t2.micro)

For deployment on t2.micro instances, the following settings are recommended:

- `GOMAXPROCS=1` to match available vCPUs
- Limited connection pool sizes
- Reduced Redis and Elasticsearch memory usage

## Build and Run

### Prerequisites

- Go 1.19 or higher
- MongoDB
- Redis
- Elasticsearch
- Kafka

### Build

```bash
# Build the application
go build -o query-service

# Build with optimizations
go build -ldflags="-s -w" -o query-service
```

### Run

```bash
# Run with default config
./query-service

# Run with custom config file
./query-service -config=/path/to/config.yaml

# Run with environment variables
MONGODB_URI=mongodb://custom-host:27017 KAFKA_BROKERS=custom-kafka:9092 ./query-service
```

### Docker

```bash
# Build Docker image
docker build -t justbrowsing/query-service .

# Run Docker container
docker run -p 8081:8081 \
  -e MONGODB_URI=mongodb://mongodb:27017/ecommerce \
  -e ELASTICSEARCH_ADDRESSES=http://elasticsearch:9200 \
  -e REDIS_ADDRESS=redis:6379 \
  -e KAFKA_BROKERS=kafka:9092 \
  -e GOMAXPROCS=1 \
  justbrowsing/query-service
```

## Monitoring

The service exposes the following monitoring endpoints:

- Health check: `/health`
- Metrics (Prometheus format): `/metrics`
- Readiness check: `/ready`

## Development Notes

### Adding a New Event Handler

1. Define the event structure in `internal/models/event.go`
2. Add a handler method in the appropriate service
3. Register the event in `internal/services/event_service.go`

### Optimizing Search Queries

For complex search operations:

1. First check Redis cache for results
2. If cache miss, try Elasticsearch for full-text search capabilities
3. If Elasticsearch is unavailable, fall back to MongoDB queries
4. Cache results for future requests

### Error Handling

The service uses a centralized error handling approach:

1. Domain-specific errors are defined in `pkg/errors`
2. The API layer converts these to appropriate HTTP responses
3. Middleware handles panic recovery and request logging

### Concurrent Request Handling

Go's concurrency model is leveraged for efficient request handling:

1. Each request is processed in a separate goroutine
2. Connection pools are used for database and cache access
3. Context is passed through the call chain for cancellation
4. Care is taken to prevent goroutine leaks