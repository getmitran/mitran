terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.0" }
  }
}

provider "aws" {
  region = var.region
}

variable "region" { default = "us-east-1" }
variable "env" { default = "prod" }
variable "vpc_cidr" { default = "10.0.0.0/16" }

# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  tags = { Name = "mitran-${var.env}" }
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  tags = { Name = "mitran-${var.env}-public" }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

# ECS Cluster
resource "aws_ecs_cluster" "mitran" {
  name = "mitran-${var.env}"
}

# ECR Repositories
resource "aws_ecr_repository" "engine" {
  name = "mitran-engine"
}

resource "aws_ecr_repository" "worker" {
  name = "mitran-worker"
}

resource "aws_ecr_repository" "dashboard" {
  name = "mitran-dashboard"
}

# IAM Role for ECS tasks
resource "aws_iam_role" "ecs_task" {
  name = "mitran-${var.env}-ecs-task"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Bedrock access for worker
resource "aws_iam_policy" "bedrock_access" {
  name = "mitran-${var.env}-bedrock"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = ["bedrock:InvokeModel", "bedrock:InvokeModelWithResponseStream"]
      Resource = "*"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "bedrock" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = aws_iam_policy.bedrock_access.arn
}

output "cluster_name" { value = aws_ecs_cluster.mitran.name }
output "ecr_engine" { value = aws_ecr_repository.engine.repository_url }
output "ecr_worker" { value = aws_ecr_repository.worker.repository_url }
output "ecr_dashboard" { value = aws_ecr_repository.dashboard.repository_url }
