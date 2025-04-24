provider "aws" {
  region = var.aws_region
}

# Configure Terraform backend for state management
terraform {
  backend "s3" {
    bucket         = "justbrowsing-terraform-state"
    key            = "terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "justbrowsing-terraform-locks"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

module "vpc" {
  source = "../modules/vpc"

  vpc_cidr_block              = var.vpc_cidr_block
  public_subnet_cidr_blocks   = var.public_subnet_cidr_blocks
  private_subnet_cidr_blocks  = var.private_subnet_cidr_blocks
  availability_zones          = var.availability_zones
  environment                 = var.environment
}

module "security" {
  source = "../modules/security"

  vpc_id                      = module.vpc.vpc_id
  environment                 = var.environment
}

module "elb" {
  source = "../modules/elb"

  vpc_id                      = module.vpc.vpc_id
  public_subnet_ids           = module.vpc.public_subnet_ids
  private_subnet_ids          = module.vpc.private_subnet_ids
  api_gateway_sg_id           = module.security.api_gateway_sg_id
  command_service_sg_id       = module.security.command_service_sg_id
  query_service_sg_id         = module.security.query_service_sg_id
  certificate_arn             = var.certificate_arn
  environment                 = var.environment
}

module "ec2" {
  source = "../modules/ec2"

  vpc_id                      = module.vpc.vpc_id
  public_subnet_ids           = module.vpc.public_subnet_ids
  private_subnet_ids          = module.vpc.private_subnet_ids
  api_gateway_sg_id           = module.security.api_gateway_sg_id
  command_service_sg_id       = module.security.command_service_sg_id
  query_service_sg_id         = module.security.query_service_sg_id
  database_sg_id              = module.security.database_sg_id
  elk_sg_id                   = module.security.elk_sg_id
  monitoring_sg_id            = module.security.monitoring_sg_id
  instance_type               = var.instance_type
  key_name                    = var.key_name
  api_gateway_target_group_arns = [module.elb.api_gateway_target_group_arn]
  command_service_target_group_arns = [module.elb.command_service_target_group_arn]
  query_service_target_group_arns = [module.elb.query_service_target_group_arn]
  environment                 = var.environment
}