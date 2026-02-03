output "instance_public_ip" {
  description = "Public IP of EC2 instance"
  value       = aws_instance.this.public_ip
}

output "instance_public_dns" {
  description = "Public DNS of EC2 instance"
  value       = aws_instance.this.public_dns
}

output "server_url" {
  description = "Full URL to access server"
  value       = "http://${aws_route53_record.this.name}.${var.domain_name}"
}

output "ssh_command" {
  description = "SSH command to connect"
  value       = "ssh -i crab-trap-key.pem ubuntu@${aws_instance.this.public_ip}"
}

output "cloudwatch_server_logs" {
  description = "CloudWatch Logs URL for server"
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#logsV2:log-groups/log-group/${replace(aws_cloudwatch_log_group.server.name, "/", "$252F")}"
}

output "cloudwatch_worker_logs" {
  description = "CloudWatch Logs URL for worker"
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#logsV2:log-groups/log-group/${replace(aws_cloudwatch_log_group.worker.name, "/", "$252F")}"
}
