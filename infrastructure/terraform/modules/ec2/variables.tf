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

variable "database_sg_id" {
  description = "ID of the Database security group"
  type        = string
}

variable "elk_sg_id" {
  description = "ID of the ELK Stack security group"
  type        = string
}

variable "monitoring_sg_id" {
  description = "ID of the Monitoring security group"
  type        = string
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

variable "key_name" {
  description = "SSH key name to use for EC2 instances"
  type        = string
}

variable "api_gateway_target_group_arns" {
  description = "ARNs of the API Gateway target groups"
  type        = list(string)
}

variable "command_service_target_group_arns" {
  description = "ARNs of the Command Service target groups"
  type        = list(string)
}

variable "query_service_target_group_arns" {
  description = "ARNs of the Query Service target groups"
  type        = list(string)
}

variable "environment" {
  description = "Environment name"
  type        = string
}