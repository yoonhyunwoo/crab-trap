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
    from_port   = 8080
    to_port     = 8080
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

data "template_file" "user_data" {
  template = file("${path.module}/user_data.sh")

  vars = {
    moltbook_api_key = var.moltbook_api_key
    moltbook_submolt = var.moltbook_submolt
    worker_interval   = var.worker_interval_minutes
    server_log_group = aws_cloudwatch_log_group.server.name
    worker_log_group = aws_cloudwatch_log_group.worker.name
    server_port      = 8080
  }
}

resource "aws_instance" "this" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.this.id]
  iam_instance_profile   = aws_iam_instance_profile.this.name

  user_data = data.template_file.user_data.rendered

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

data "aws_vpc" "default" {
  default = true
}
