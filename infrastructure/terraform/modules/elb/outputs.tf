output "api_gateway_dns_name" {
  description = "DNS name of the public ALB"
  value       = aws_lb.public.dns_name
}

output "command_service_dns_name" {
  description = "DNS name of the Command Service ALB"
  value       = aws_lb.command_service.dns_name
}

output "query_service_dns_name" {
  description = "DNS name of the Query Service ALB"
  value       = aws_lb.query_service.dns_name
}

output "api_gateway_target_group_arn" {
  description = "ARN of the API Gateway target group"
  value       = aws_lb_target_group.api_gateway.arn
}

output "command_service_target_group_arn" {
  description = "ARN of the Command Service target group"
  value       = aws_lb_target_group.command_service.arn
}

output "query_service_target_group_arn" {
  description = "ARN of the Query Service target group"
  value       = aws_lb_target_group.query_service.arn
}