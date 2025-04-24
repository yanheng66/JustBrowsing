# Public-facing Application Load Balancer
resource "aws_lb" "public" {
  name               = "${var.environment}-public-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.api_gateway_sg_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = false

  tags = {
    Name        = "${var.environment}-public-alb"
    Environment = var.environment
  }
}

# API Gateway Target Group
resource "aws_lb_target_group" "api_gateway" {
  name     = "${var.environment}-api-gateway-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = var.vpc_id

  health_check {
    path                = "/health"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200"
  }

  tags = {
    Name        = "${var.environment}-api-gateway-tg"
    Environment = var.environment
  }
}

# Public ALB Listeners
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.public.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_gateway.arn
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.public.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = var.certificate_arn  # Optional - needs to be provided

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_gateway.arn
  }
}

# Internal Command Service Load Balancer
resource "aws_lb" "command_service" {
  name               = "${var.environment}-command-service-alb"
  internal           = true
  load_balancer_type = "application"
  security_groups    = [var.command_service_sg_id]
  subnets            = var.private_subnet_ids

  enable_deletion_protection = false

  tags = {
    Name        = "${var.environment}-command-service-alb"
    Environment = var.environment
  }
}

# Command Service Target Group
resource "aws_lb_target_group" "command_service" {
  name     = "${var.environment}-command-service-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = var.vpc_id

  health_check {
    path                = "/api/commands/actuator/health"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200"
  }

  tags = {
    Name        = "${var.environment}-command-service-tg"
    Environment = var.environment
  }
}

# Command Service Listener
resource "aws_lb_listener" "command_service" {
  load_balancer_arn = aws_lb.command_service.arn
  port              = 8080
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.command_service.arn
  }
}

# Internal Query Service Load Balancer
resource "aws_lb" "query_service" {
  name               = "${var.environment}-query-service-alb"
  internal           = true
  load_balancer_type = "application"
  security_groups    = [var.query_service_sg_id]
  subnets            = var.private_subnet_ids

  enable_deletion_protection = false

  tags = {
    Name        = "${var.environment}-query-service-alb"
    Environment = var.environment
  }
}

# Query Service Target Group
resource "aws_lb_target_group" "query_service" {
  name     = "${var.environment}-query-service-tg"
  port     = 8081
  protocol = "HTTP"
  vpc_id   = var.vpc_id

  health_check {
    path                = "/health"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200"
  }

  tags = {
    Name        = "${var.environment}-query-service-tg"
    Environment = var.environment
  }
}

# Query Service Listener
resource "aws_lb_listener" "query_service" {
  load_balancer_arn = aws_lb.query_service.arn
  port              = 8081
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.query_service.arn
  }
}