output "vpc_id" {
  description = "ID of the created VPC"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

output "api_gateway_public_dns" {
  description = "Public DNS of the API Gateway load balancer"
  value       = module.elb.api_gateway_dns_name
}

output "api_gateway_instances" {
  description = "API Gateway EC2 instances"
  value       = module.ec2.api_gateway_instances
}

output "command_service_instances" {
  description = "Command Service EC2 instances"
  value       = module.ec2.command_service_instances
}

output "query_service_instances" {
  description = "Query Service EC2 instances"
  value       = module.ec2.query_service_instances
}

output "database_instances" {
  description = "Database EC2 instances"
  value       = module.ec2.database_instances
}

output "elk_instances" {
  description = "ELK Stack EC2 instances"
  value       = module.ec2.elk_instances
}

output "monitoring_instances" {
  description = "Monitoring EC2 instances"
  value       = module.ec2.monitoring_instances
}