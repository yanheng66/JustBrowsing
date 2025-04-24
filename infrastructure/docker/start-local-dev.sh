#!/bin/bash

# Create database initialization scripts directory
mkdir -p ./databases/postgres
mkdir -p ./databases/mongodb
mkdir -p ./databases/redis
mkdir -p ./services/nginx/logs

# Create Redis configuration
cat > ./databases/redis/redis.conf << EOL
# Redis configuration for local development
maxmemory 256mb
maxmemory-policy allkeys-lru
EOL

# Create PostgreSQL initialization script
cat > ./databases/postgres/init.sql << EOL
-- PostgreSQL initialization script
CREATE DATABASE ecommerce;
\c ecommerce;

-- Create tables based on the database schema
-- Products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ProductTags table
CREATE TABLE product_tags (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    tag_id INTEGER REFERENCES tags(id),
    tag_value VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (product_id, tag_id)
);

-- Inventory table
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 0,
    last_replenishment_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_quantity CHECK (quantity >= 0)
);

-- Orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- OrderItems table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    CONSTRAINT positive_quantity CHECK (quantity > 0)
);

-- Outbox table
CREATE TABLE outbox_events (
    id SERIAL PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_inventory_product_id ON inventory(product_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_outbox_events_processed_created_at ON outbox_events(processed, created_at);
CREATE INDEX idx_product_tags_product_id_tag_id ON product_tags(product_id, tag_id);
CREATE INDEX idx_tags_name ON tags(name);
EOL

# Create MongoDB initialization script
cat > ./databases/mongodb/init.js << EOL
// Create ecommerce database and collections
db = db.getSiblingDB('ecommerce');

// Products collection
db.createCollection('products');
db.products.createIndex({ "productId": 1 }, { unique: true });
db.products.createIndex({ "sku": 1 }, { unique: true });
db.products.createIndex({ "tags.name": 1, "tags.value": 1 });

// Orders collection
db.createCollection('orders');
db.orders.createIndex({ "orderId": 1 }, { unique: true });
db.orders.createIndex({ "orderNumber": 1 }, { unique: true });
EOL

# Ensure Docker network exists
docker network create ecommerce-network 2>/dev/null || true

# Start core services
echo "Starting core services (databases, Kafka, Elasticsearch)..."
docker-compose up -d postgres mongodb redis zookeeper kafka schema-registry elasticsearch

# Wait for core services to be ready
echo "Waiting for core services to be ready..."
sleep 20

# Start application services
echo "Starting application services..."
docker-compose up -d command-service query-service api-gateway

# Start monitoring services
echo "Starting monitoring services..."
docker-compose -f docker-compose.monitoring.yml up -d

echo "All services started. The application is available at: http://localhost"
echo "Monitoring dashboard is available at: http://localhost:3000 (admin/admin)"
echo "Kibana is available at: http://localhost:5601"
echo "Jaeger UI is available at: http://localhost:16686"