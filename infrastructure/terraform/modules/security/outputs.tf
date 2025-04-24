output "api_gateway_sg_id" {
  description = "ID of the API Gateway security group"
  value       = aws_security_group.api_gateway.id
}

output "command_service_sg_id" {
  description = "ID of the Command Service security group"
  value       = aws_security_group.command_service.id
}

output "query_service_sg_id" {
  description = "ID of the Query Service security group"
  value       = aws_security_group.query_service.id
}

output "database_sg_id" {
  description = "ID of the Database security group"
  value       = aws_security_group.database.id
}

output "elk_sg_id" {
  description = "ID of the ELK Stack security group"
  value       = aws_security_group.elk.id
}

output "monitoring_sg_id" {
  description = "ID of the Monitoring security group"
  value       = aws_security_group.monitoring.id
}