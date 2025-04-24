output "api_gateway_instances" {
  description = "API Gateway EC2 instances"
  value = [
    {
      id        = aws_instance.api_gateway_monitoring_1.id
      name      = aws_instance.api_gateway_monitoring_1.tags.Name
      public_ip = aws_instance.api_gateway_monitoring_1.public_ip
      private_ip = aws_instance.api_gateway_monitoring_1.private_ip
      subnet_id = aws_instance.api_gateway_monitoring_1.subnet_id
      role      = "api-gateway,monitoring"
    },
    {
      id        = aws_instance.api_gateway_2.id
      name      = aws_instance.api_gateway_2.tags.Name
      public_ip = aws_instance.api_gateway_2.public_ip
      private_ip = aws_instance.api_gateway_2.private_ip
      subnet_id = aws_instance.api_gateway_2.subnet_id
      role      = "api-gateway"
    }
  ]
}

output "command_service_instances" {
  description = "Command Service EC2 instances"
  value = [
    {
      id        = aws_instance.command_service_kafka_1.id
      name      = aws_instance.command_service_kafka_1.tags.Name
      private_ip = aws_instance.command_service_kafka_1.private_ip
      subnet_id = aws_instance.command_service_kafka_1.subnet_id
      role      = "command-service,kafka"
    },
    {
      id        = aws_instance.command_service_2.id
      name      = aws_instance.command_service_2.tags.Name
      private_ip = aws_instance.command_service_2.private_ip
      subnet_id = aws_instance.command_service_2.subnet_id
      role      = "command-service"
    }
  ]
}

output "query_service_instances" {
  description = "Query Service EC2 instances"
  value = [
    {
      id        = aws_instance.query_service_redis_1.id
      name      = aws_instance.query_service_redis_1.tags.Name
      private_ip = aws_instance.query_service_redis_1.private_ip
      subnet_id = aws_instance.query_service_redis_1.subnet_id
      role      = "query-service,redis"
    },
    {
      id        = aws_instance.query_service_2.id
      name      = aws_instance.query_service_2.tags.Name
      private_ip = aws_instance.query_service_2.private_ip
      subnet_id = aws_instance.query_service_2.subnet_id
      role      = "query-service"
    }
  ]
}

output "database_instances" {
  description = "Database EC2 instances"
  value = [
    {
      id        = aws_instance.postgres_elk.id
      name      = aws_instance.postgres_elk.tags.Name
      private_ip = aws_instance.postgres_elk.private_ip
      subnet_id = aws_instance.postgres_elk.subnet_id
      role      = "database-postgres,elk-stack"
    },
    {
      id        = aws_instance.mongodb_elasticsearch.id
      name      = aws_instance.mongodb_elasticsearch.tags.Name
      private_ip = aws_instance.mongodb_elasticsearch.private_ip
      subnet_id = aws_instance.mongodb_elasticsearch.subnet_id
      role      = "database-mongodb,elasticsearch"
    }
  ]
}