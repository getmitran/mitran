variable "region" {
  default = "us-east-1"
}

variable "cluster_name" {
  default = "mitran-cluster"
}

variable "image_tag" {
  default = "latest"
}

provider "aws" {
  region = var.region
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_cloudwatch_log_group" "mitran" {
  name              = "/ecs/mitran"
  retention_in_days = 14
}

resource "aws_security_group" "mitran" {
  name   = "mitran-sg"
  vpc_id = data.aws_vpc.default.id

  ingress {
    from_port   = 7780
    to_port     = 7780
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8888
    to_port     = 8888
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_cluster" "mitran" {
  name = var.cluster_name
}

resource "aws_ecs_task_definition" "mitran" {
  family                   = "mitran"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.ecs_execution.arn

  container_definitions = jsonencode([
    {
      name      = "mitran-engine"
      image     = "mitran/engine:${var.image_tag}"
      essential = true
      portMappings = [{ containerPort = 7780, protocol = "tcp" }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.mitran.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "engine"
        }
      }
    },
    {
      name      = "mitran-worker"
      image     = "mitran/worker:${var.image_tag}"
      essential = true
      portMappings = [{ containerPort = 8888, protocol = "tcp" }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.mitran.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "worker"
        }
      }
    }
  ])
}

resource "aws_iam_role" "ecs_execution" {
  name = "mitran-ecs-execution"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_execution" {
  role       = aws_iam_role.ecs_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_lb" "mitran" {
  name               = "mitran-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.mitran.id]
  subnets            = data.aws_subnets.default.ids
}

resource "aws_lb_target_group" "mitran" {
  name        = "mitran-tg"
  port        = 7780
  protocol    = "HTTP"
  vpc_id      = data.aws_vpc.default.id
  target_type = "ip"

  health_check {
    path = "/healthz"
  }
}

resource "aws_lb_listener" "mitran" {
  load_balancer_arn = aws_lb.mitran.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.mitran.arn
  }
}

resource "aws_ecs_service" "mitran" {
  name            = "mitran-service"
  cluster         = aws_ecs_cluster.mitran.id
  task_definition = aws_ecs_task_definition.mitran.arn
  desired_count   = 2
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = data.aws_subnets.default.ids
    security_groups = [aws_security_group.mitran.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.mitran.arn
    container_name   = "mitran-engine"
    container_port   = 7780
  }
}

output "alb_dns_name" {
  value = aws_lb.mitran.dns_name
}
