variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "public_subnet_ids" {
  description = "IDs of the public subnets"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "IDs of the private subnets"
  type        = list(string)
}

variable "api_gateway_sg_id" {
  description = "ID of the API Gateway security group"
  type        = string
}

variable "command_service_sg_id" {
  description = "ID of the Command Service security group"
  type        = string
}

variable "query_service_sg_id" {
  description = "ID of the Query Service security group"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "certificate_arn" {
  description = "ARN of the SSL certificate for HTTPS"
  type        = string
  default     = null
}