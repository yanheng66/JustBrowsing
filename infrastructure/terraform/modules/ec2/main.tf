# AMIs
locals {
  amazon_linux_ami = "ami-0c55b159cbfafe1f0"  # Amazon Linux 2 for us-east-1
}

# Instance 1: API Gateway, Prometheus, Grafana (Public Subnet)
resource "aws_instance" "api_gateway_monitoring_1" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.api_gateway_sg_id, var.monitoring_sg_id]
  subnet_id              = var.public_subnet_ids[0]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 15
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    systemctl start docker
    systemctl enable docker
    
    # Install Docker Compose
    curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    # Set up API Gateway
    mkdir -p /opt/api-gateway/nginx
    cat > /opt/api-gateway/nginx/nginx.conf << 'EOL'
    ${file("../../../../api-gateway/nginx/nginx.conf")}
    EOL
    
    cat > /opt/api-gateway/docker-compose.yml << 'EOL'
    ${file("../../../../api-gateway/docker-compose.yml")}
    EOL
    
    # Set up Prometheus and Grafana
    mkdir -p /opt/monitoring/{prometheus,grafana}
    
    cat > /opt/monitoring/docker-compose.yml << 'EOL'
    version: '3'
    services:
      prometheus:
        image: prom/prometheus:v2.45.0
        ports:
          - "9090:9090"
        volumes:
          - /opt/monitoring/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
          - /opt/monitoring/prometheus/data:/prometheus
        command:
          - '--config.file=/etc/prometheus/prometheus.yml'
          - '--storage.tsdb.path=/prometheus'
          - '--web.console.libraries=/etc/prometheus/console_libraries'
          - '--web.console.templates=/etc/prometheus/consoles'
          - '--web.enable-lifecycle'
          - '--storage.tsdb.retention.time=15d'
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512

      grafana:
        image: grafana/grafana:10.0.3
        ports:
          - "3000:3000"
        environment:
          - GF_SECURITY_ADMIN_USER=admin
          - GF_SECURITY_ADMIN_PASSWORD=admin
          - GF_USERS_ALLOW_SIGN_UP=false
        volumes:
          - /opt/monitoring/grafana:/var/lib/grafana
        depends_on:
          - prometheus
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
    EOL
    
    # Create Prometheus config
    cat > /opt/monitoring/prometheus/prometheus.yml << 'EOL'
    global:
      scrape_interval: 15s
      evaluation_interval: 15s

    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
          - targets: ['localhost:9090']

      - job_name: 'command-service'
        metrics_path: '/api/commands/actuator/prometheus'
        static_configs:
          - targets: ['command-service-1:8080', 'command-service-2:8080']

      - job_name: 'query-service'
        metrics_path: '/metrics'
        static_configs:
          - targets: ['query-service-1:8081', 'query-service-2:8081']
    EOL
    
    # Start services
    cd /opt/api-gateway
    docker-compose up -d
    
    cd /opt/monitoring
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-api-gateway-monitoring-1"
    Environment = var.environment
    Role        = "api-gateway,monitoring"
  }
}

# Instance 2: API Gateway Redundancy (Public Subnet)
resource "aws_instance" "api_gateway_2" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.api_gateway_sg_id]
  subnet_id              = var.public_subnet_ids[1]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 10
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    systemctl start docker
    systemctl enable docker
    
    # Install Docker Compose
    curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    # Set up API Gateway
    mkdir -p /opt/api-gateway/nginx
    cat > /opt/api-gateway/nginx/nginx.conf << 'EOL'
    ${file("../../../../api-gateway/nginx/nginx.conf")}
    EOL
    
    cat > /opt/api-gateway/docker-compose.yml << 'EOL'
    ${file("../../../../api-gateway/docker-compose.yml")}
    EOL
    
    # Start the API Gateway
    cd /opt/api-gateway
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-api-gateway-2"
    Environment = var.environment
    Role        = "api-gateway"
  }
}

# Instance 3: Command Service, Kafka (Private Subnet)
resource "aws_instance" "command_service_kafka_1" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.command_service_sg_id]
  subnet_id              = var.private_subnet_ids[0]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 15
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    amazon-linux-extras install -y java-openjdk11
    systemctl start docker
    systemctl enable docker
    
    # Set JVM options
    echo 'JAVA_OPTS="-Xms256m -Xmx512m -XX:+UseCompressedOops -Xss256k"' > /etc/environment
    
    # Set up Kafka
    mkdir -p /opt/kafka/data
    
    cat > /opt/kafka/docker-compose.yml << 'EOL'
    version: '3'
    services:
      zookeeper:
        image: confluentinc/cp-zookeeper:7.4.0
        ports:
          - "2181:2181"
        environment:
          ZOOKEEPER_CLIENT_PORT: 2181
          ZOOKEEPER_TICK_TIME: 2000
        volumes:
          - /opt/kafka/data/zookeeper:/var/lib/zookeeper/data
        restart: unless-stopped
        mem_limit: 256m
        cpu_shares: 512

      kafka:
        image: confluentinc/cp-kafka:7.4.0
        depends_on:
          - zookeeper
        ports:
          - "9092:9092"
        environment:
          KAFKA_BROKER_ID: 1
          KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
          KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:9092
          KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
          KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
          KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
          KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
          KAFKA_HEAP_OPTS: "-Xms256m -Xmx512m"
        volumes:
          - /opt/kafka/data/kafka:/var/lib/kafka/data
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
        
      schema-registry:
        image: confluentinc/cp-schema-registry:7.4.0
        depends_on:
          - kafka
        ports:
          - "8081:8081"
        environment:
          SCHEMA_REGISTRY_HOST_NAME: schema-registry
          SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS: kafka:9092
          SCHEMA_REGISTRY_LISTENERS: http://0.0.0.0:8081
        restart: unless-stopped
        mem_limit: 256m
        cpu_shares: 512
    EOL
    
    # Set up Command Service
    mkdir -p /opt/command-service
    
    cat > /opt/command-service/docker-compose.yml << 'EOL'
    version: '3'
    services:
      command-service:
        image: justbrowsing/command-service:latest
        ports:
          - "8080:8080"
        environment:
          - SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/ecommerce
          - SPRING_DATASOURCE_USERNAME=postgres
          - SPRING_DATASOURCE_PASSWORD=postgres
          - SPRING_KAFKA_BOOTSTRAP_SERVERS=kafka:9092
          - SPRING_KAFKA_PRODUCER_PROPERTIES_SCHEMA_REGISTRY_URL=http://schema-registry:8081
          - JAVA_OPTS=-Xms256m -Xmx512m -XX:+UseCompressedOops -Xss256k
        healthcheck:
          test: ["CMD", "curl", "-f", "http://localhost:8080/api/commands/actuator/health"]
          interval: 30s
          timeout: 10s
          retries: 3
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
    EOL
    
    # Start Kafka and Schema Registry
    cd /opt/kafka
    docker-compose up -d
    
    # Wait for Kafka to be ready
    sleep 20
    
    # Start Command Service
    cd /opt/command-service
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-command-service-kafka-1"
    Environment = var.environment
    Role        = "command-service,kafka"
  }
}

# Instance 4: Command Service Redundancy (Private Subnet)
resource "aws_instance" "command_service_2" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.command_service_sg_id]
  subnet_id              = var.private_subnet_ids[1]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 12
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    amazon-linux-extras install -y java-openjdk11
    systemctl start docker
    systemctl enable docker
    
    # Set JVM options
    echo 'JAVA_OPTS="-Xms256m -Xmx512m -XX:+UseCompressedOops -Xss256k"' > /etc/environment
    
    # Set up Command Service
    mkdir -p /opt/command-service
    
    cat > /opt/command-service/docker-compose.yml << 'EOL'
    version: '3'
    services:
      command-service:
        image: justbrowsing/command-service:latest
        ports:
          - "8080:8080"
        environment:
          - SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/ecommerce
          - SPRING_DATASOURCE_USERNAME=postgres
          - SPRING_DATASOURCE_PASSWORD=postgres
          - SPRING_KAFKA_BOOTSTRAP_SERVERS=kafka:9092
          - SPRING_KAFKA_PRODUCER_PROPERTIES_SCHEMA_REGISTRY_URL=http://schema-registry:8081
          - JAVA_OPTS=-Xms256m -Xmx512m -XX:+UseCompressedOops -Xss256k
        healthcheck:
          test: ["CMD", "curl", "-f", "http://localhost:8080/api/commands/actuator/health"]
          interval: 30s
          timeout: 10s
          retries: 3
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
    EOL
    
    # Start Command Service
    cd /opt/command-service
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-command-service-2"
    Environment = var.environment
    Role        = "command-service"
  }
}

# Instance 5: Query Service, Redis (Private Subnet)
resource "aws_instance" "query_service_redis_1" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.query_service_sg_id]
  subnet_id              = var.private_subnet_ids[2]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 15
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    systemctl start docker
    systemctl enable docker
    
    # Install Docker Compose
    curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    # Set up Redis
    mkdir -p /opt/redis/data
    
    cat > /opt/redis/redis.conf << 'EOL'
    # Redis configuration for resource-constrained environment
    maxmemory 256mb
    maxmemory-policy allkeys-lru
    EOL
    
    cat > /opt/redis/docker-compose.yml << 'EOL'
    version: '3'
    services:
      redis:
        image: redis:7-alpine
        ports:
          - "6379:6379"
        volumes:
          - /opt/redis/data:/data
          - /opt/redis/redis.conf:/usr/local/etc/redis/redis.conf
        command: ["redis-server", "/usr/local/etc/redis/redis.conf"]
        restart: unless-stopped
        mem_limit: 256m
        cpu_shares: 512
    EOL
    
    # Set up Query Service
    mkdir -p /opt/query-service
    
    cat > /opt/query-service/docker-compose.yml << 'EOL'
    version: '3'
    services:
      query-service:
        image: justbrowsing/query-service:latest
        ports:
          - "8081:8081"
        environment:
          - MONGODB_URI=mongodb://mongodb:27017/ecommerce
          - ELASTICSEARCH_ADDRESSES=http://elasticsearch:9200
          - REDIS_ADDRESS=redis:6379
          - KAFKA_BROKERS=kafka:9092
          - GOMAXPROCS=1
        healthcheck:
          test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
          interval: 30s
          timeout: 10s
          retries: 3
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
    EOL
    
    # Start Redis
    cd /opt/redis
    docker-compose up -d
    
    # Start Query Service
    cd /opt/query-service
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-query-service-redis-1"
    Environment = var.environment
    Role        = "query-service,redis"
  }
}

# Instance 6: Query Service Redundancy (Private Subnet)
resource "aws_instance" "query_service_2" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.query_service_sg_id]
  subnet_id              = var.private_subnet_ids[3]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 12
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    systemctl start docker
    systemctl enable docker
    
    # Install Docker Compose
    curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    # Set up Query Service
    mkdir -p /opt/query-service
    
    cat > /opt/query-service/docker-compose.yml << 'EOL'
    version: '3'
    services:
      query-service:
        image: justbrowsing/query-service:latest
        ports:
          - "8081:8081"
        environment:
          - MONGODB_URI=mongodb://mongodb:27017/ecommerce
          - ELASTICSEARCH_ADDRESSES=http://elasticsearch:9200
          - REDIS_ADDRESS=redis:6379
          - KAFKA_BROKERS=kafka:9092
          - GOMAXPROCS=1
        healthcheck:
          test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
          interval: 30s
          timeout: 10s
          retries: 3
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
    EOL
    
    # Start Query Service
    cd /opt/query-service
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-query-service-2"
    Environment = var.environment
    Role        = "query-service"
  }
}

# Instance 7: PostgreSQL, ELK Stack (Private Subnet)
resource "aws_instance" "postgres_elk" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.database_sg_id, var.elk_sg_id]
  subnet_id              = var.private_subnet_ids[0]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 20
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    systemctl start docker
    systemctl enable docker
    
    # Install Docker Compose
    curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    # Set up PostgreSQL
    mkdir -p /opt/postgres/data
    
    cat > /opt/postgres/init.sql << 'EOL'
    -- Add PostgreSQL initialization script here
    -- Create tables, indexes, etc.
    ${file("../../../../infrastructure/docker/databases/postgres/init.sql")}
    EOL
    
    cat > /opt/postgres/docker-compose.yml << 'EOL'
    version: '3'
    services:
      postgres:
        image: postgres:14-alpine
        ports:
          - "5432:5432"
        environment:
          POSTGRES_DB: ecommerce
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        volumes:
          - /opt/postgres/data:/var/lib/postgresql/data
          - /opt/postgres/init.sql:/docker-entrypoint-initdb.d/init.sql
        restart: unless-stopped
        command: ["postgres", "-c", "max_connections=50", "-c", "shared_buffers=128MB"]
        mem_limit: 256m
        cpu_shares: 512
    EOL
    
    # Set up ELK Stack
    mkdir -p /opt/elk/{elasticsearch,logstash,kibana,logstash-pipeline}
    
    cat > /opt/elk/logstash-pipeline/logstash.conf << 'EOL'
    input {
      beats {
        port => 5044
      }
      tcp {
        port => 5000
        codec => json
      }
    }
    
    filter {
      if [type] == "container" {
        # Parse container logs
        if [log] =~ "^{" {
          json {
            source => "log"
          }
        }
      }
    }
    
    output {
      elasticsearch {
        hosts => ["elasticsearch:9200"]
        index => "filebeat-%{+YYYY.MM.dd}"
      }
    }
    EOL
    
    cat > /opt/elk/docker-compose.yml << 'EOL'
    version: '3'
    services:
      elasticsearch:
        image: docker.elastic.co/elasticsearch/elasticsearch:8.7.1
        ports:
          - "9200:9200"
        environment:
          - discovery.type=single-node
          - bootstrap.memory_lock=true
          - "ES_JAVA_OPTS=-Xms256m -Xmx512m"
          - xpack.security.enabled=false
          - xpack.ml.enabled=false
        volumes:
          - /opt/elk/elasticsearch:/usr/share/elasticsearch/data
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
        ulimits:
          memlock:
            soft: -1
            hard: -1
    
      logstash:
        image: docker.elastic.co/logstash/logstash:8.7.1
        ports:
          - "5044:5044"
          - "5000:5000"
        environment:
          - "LS_JAVA_OPTS=-Xms128m -Xmx256m"
        volumes:
          - /opt/elk/logstash-pipeline:/usr/share/logstash/pipeline
        depends_on:
          - elasticsearch
        restart: unless-stopped
        mem_limit: 256m
        cpu_shares: 512
    
      kibana:
        image: docker.elastic.co/kibana/kibana:8.7.1
        ports:
          - "5601:5601"
        environment:
          - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
        depends_on:
          - elasticsearch
        restart: unless-stopped
        mem_limit: 256m
        cpu_shares: 512
    EOL
    
    # Start PostgreSQL
    cd /opt/postgres
    docker-compose up -d
    
    # Start ELK Stack
    cd /opt/elk
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-postgres-elk"
    Environment = var.environment
    Role        = "database,logging"
  }
}

# Instance 8: MongoDB, Elasticsearch (Private Subnet)
resource "aws_instance" "mongodb_elasticsearch" {
  ami                    = local.amazon_linux_ami
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [var.database_sg_id]
  subnet_id              = var.private_subnet_ids[1]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 20
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    systemctl start docker
    systemctl enable docker
    
    # Install Docker Compose
    curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    # Set up MongoDB
    mkdir -p /opt/mongodb/data
    
    cat > /opt/mongodb/init.js << 'EOL'
    // MongoDB initialization script
    ${file("../../../../infrastructure/docker/databases/mongodb/init.js")}
    EOL
    
    cat > /opt/mongodb/docker-compose.yml << 'EOL'
    version: '3'
    services:
      mongodb:
        image: mongo:6.0
        ports:
          - "27017:27017"
        environment:
          MONGO_INITDB_DATABASE: ecommerce
        volumes:
          - /opt/mongodb/data:/data/db
          - /opt/mongodb/init.js:/docker-entrypoint-initdb.d/init.js
        restart: unless-stopped
        command: ["--wiredTigerCacheSizeGB", "0.25"]
        mem_limit: 256m
        cpu_shares: 512
    EOL
    
    # Set up Elasticsearch for search
    mkdir -p /opt/elasticsearch/data
    
    cat > /opt/elasticsearch/docker-compose.yml << 'EOL'
    version: '3'
    services:
      elasticsearch:
        image: docker.elastic.co/elasticsearch/elasticsearch:8.7.1
        ports:
          - "9200:9200"
          - "9300:9300"
        environment:
          - discovery.type=single-node
          - bootstrap.memory_lock=true
          - "ES_JAVA_OPTS=-Xms256m -Xmx512m"
          - xpack.security.enabled=false
          - xpack.ml.enabled=false
        volumes:
          - /opt/elasticsearch/data:/usr/share/elasticsearch/data
        restart: unless-stopped
        mem_limit: 512m
        cpu_shares: 512
        ulimits:
          memlock:
            soft: -1
            hard: -1
    EOL
    
    # Start MongoDB
    cd /opt/mongodb
    docker-compose up -d
    
    # Start Elasticsearch
    cd /opt/elasticsearch
    docker-compose up -d
  EOF

  tags = {
    Name        = "${var.environment}-mongodb-elasticsearch"
    Environment = var.environment
    Role        = "database,search"
  }
}

# Register API Gateway instances with target group
resource "aws_lb_target_group_attachment" "api_gateway_1" {
  target_group_arn = var.api_gateway_target_group_arns[0]
  target_id        = aws_instance.api_gateway_monitoring_1.id
  port             = 80
}

resource "aws_lb_target_group_attachment" "api_gateway_2" {
  target_group_arn = var.api_gateway_target_group_arns[0]
  target_id        = aws_instance.api_gateway_2.id
  port             = 80
}

# Register Command Service instances with target group
resource "aws_lb_target_group_attachment" "command_service_1" {
  target_group_arn = var.command_service_target_group_arns[0]
  target_id        = aws_instance.command_service_kafka_1.id
  port             = 8080
}

resource "aws_lb_target_group_attachment" "command_service_2" {
  target_group_arn = var.command_service_target_group_arns[0]
  target_id        = aws_instance.command_service_2.id
  port             = 8080
}

# Register Query Service instances with target group
resource "aws_lb_target_group_attachment" "query_service_1" {
  target_group_arn = var.query_service_target_group_arns[0]
  target_id        = aws_instance.query_service_redis_1.id
  port             = 8081
}

resource "aws_lb_target_group_attachment" "query_service_2" {
  target_group_arn = var.query_service_target_group_arns[0]
  target_id        = aws_instance.query_service_2.id
  port             = 8081
}