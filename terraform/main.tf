data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_iam_role" "ec2_role" {
  name = "crab-trap-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "cloudwatch_logs" {
  name = "cloudwatch-logs-policy"
  role = aws_iam_role.ec2_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "ecr_access" {
  name = "ecr-access-policy"
  role = aws_iam_role.ec2_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability"
        ]
        Resource = aws_ecr_repository.this.arn
      },
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = aws_secretsmanager_secret.moltbook.arn
      }
    ]
  })
}

resource "aws_iam_instance_profile" "this" {
  name = "crab-trap-instance-profile"
  role = aws_iam_role.ec2_role.name
}

resource "aws_security_group" "this" {
  name        = "crab-trap-sg"
  description = "Security group for crab-trap"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.ssh_allowed_cidr
  }

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "All outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "crab-trap-sg"
  }
}

resource "aws_cloudwatch_log_group" "server" {
  name              = "/crab-trap/server"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "worker" {
  name              = "/crab-trap/worker"
  retention_in_days = 7
}

resource "aws_key_pair" "this" {
  key_name   = "crab-trap-key"
  public_key = file("${path.module}/crab-trap-key.pub")
}

locals {
  user_data = templatefile("${path.module}/user_data.sh", {
    moltbook_secret_name = aws_secretsmanager_secret.moltbook.name
    moltbook_submolt    = var.moltbook_submolt
    worker_interval      = var.worker_interval_minutes
    server_log_group    = aws_cloudwatch_log_group.server.name
    worker_log_group    = aws_cloudwatch_log_group.worker.name
    server_port         = 8080
    ecr_repository_url  = aws_ecr_repository.this.repository_url
    image_tag           = var.image_tag
    subdomain           = var.subdomain
    domain_name         = var.domain_name
  })
}

resource "aws_instance" "this" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.this.id]
  iam_instance_profile   = aws_iam_instance_profile.this.name
  key_name              = aws_key_pair.this.key_name

  user_data = local.user_data

  root_block_device {
    volume_type           = "gp3"
    volume_size           = 20
    delete_on_termination = true
  }

  tags = {
    Name        = "crab-trap"
    Environment = "production"
  }
}

resource "aws_route53_record" "this" {
  zone_id = var.hosted_zone_id
  name    = var.subdomain
  type    = "A"
  ttl     = 300
  records = [aws_instance.this.public_ip]
}

resource "aws_ecr_repository" "this" {
  name                 = "crab-trap"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "crab-trap"
  }
}

resource "aws_ecr_lifecycle_policy" "this" {
  repository = aws_ecr_repository.this.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 10 images"
        selection = {
          tagStatus     = "any"
          countType     = "imageCountMoreThan"
          countNumber   = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

resource "aws_secretsmanager_secret" "moltbook" {
  name = "crab-trap/moltbook-api-key"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "this" {
  secret_id     = aws_secretsmanager_secret.moltbook.id
  secret_string = var.moltbook_api_key
}

data "aws_vpc" "default" {
  default = true
}
