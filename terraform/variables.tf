variable "aws_region" {
  description = "AWS region"
  default     = "ap-northeast-2"
}

variable "instance_type" {
  description = "EC2 instance type"
  default     = "t3.small"
}

variable "hosted_zone_id" {
  description = "Route53 hosted zone ID"
  default     = "Z02084782RE0P3KE726F9"
}

variable "domain_name" {
  description = "Root domain name"
  default     = "thumbgo.kr"
}

variable "subdomain" {
  description = "Subdomain for injector"
  default     = "injector"
}

variable "ssh_allowed_cidr" {
  description = "CIDR block allowed to SSH"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "moltbook_api_key" {
  description = "Moltbook API key (sensitive)"
  type        = string
  sensitive   = true
}

variable "moltbook_submolt" {
  description = "Moltbook submolt name"
  default     = "general"
}

variable "worker_interval_minutes" {
  description = "Worker execution interval in minutes"
  default     = 60
}
