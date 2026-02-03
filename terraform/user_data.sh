#!/bin/bash
set -e

echo "Initializing Crab Trap..."

# Update system
apt-get update -y
apt-get install -y docker.io curl git awscli

# Start Docker
systemctl start docker
systemctl enable docker

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Create application directory
mkdir -p /opt/crab-trap
cd /opt/crab-trap


# Get API key from Secrets Manager
MOLTBOOK_API_KEY=$(aws secretsmanager get-secret-value --secret-id "${moltbook_secret_name}" --query SecretString --output text --region ap-northeast-2)

# Create config.yaml
cat > config.yaml <<'CONFIG_EOF'
server:
  port: 8080
  log_dir: "./logs"

worker:
  moltbook_api_key: "PLACEHOLDER"
  server_url: "http://localhost:8080"
  submolt: "${moltbook_submolt}"
  interval_minutes: ${worker_interval}
  os_detection: true

logging:
  level: "info"
  save_requests: true
CONFIG_EOF

# Replace placeholder with actual API key
sed -i "s/PLACEHOLDER/$MOLTBOOK_API_KEY/g" config.yaml

# Login to ECR
AWS_REGION=$(curl -s http://169.254.169.254/latest/meta-data/placement/region)
ECR_REGISTRY=$(echo "${ecr_repository_url}" | cut -d/ -f1)
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# Create docker-compose.yml
cat > docker-compose.yml <<'COMPOSE_EOF'
version: '3.8'

services:
  caddy:
    image: caddy:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config
    restart: unless-stopped
    logging:
      driver: "awslogs"
      options:
        awslogs-region: ap-northeast-2
        awslogs-group: "${server_log_group}"

  server:
    image: ${ecr_repository_url}:${image_tag}
    command: ["./server"]
    ports:
      - "127.0.0.1:8080:8080"
    environment:
      - PORT=8080
      - LOG_DIR=/logs
    volumes:
      - ./logs:/logs
    restart: unless-stopped
    logging:
      driver: "awslogs"
      options:
        awslogs-region: ap-northeast-2
        awslogs-group: "${server_log_group}"

  worker:
    image: ${ecr_repository_url}:${image_tag}
    command: ["./worker"]
    environment:
      - CONFIG=/app/config.yaml
      - LOG_DIR=/logs
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./logs:/logs
    depends_on:
      - server
    restart: unless-stopped
    logging:
      driver: "awslogs"
      options:
        awslogs-region: ap-northeast-2
        awslogs-group: "${worker_log_group}"

volumes:
  logs:
  caddy_data:
  caddy_config:
COMPOSE_EOF

# Create Caddyfile
cat > Caddyfile <<'CADDY_EOF'
${subdomain}.${domain_name} {
    reverse_proxy server:8080
}
CADDY_EOF

# Create log directory
mkdir -p /opt/crab-trap/logs

echo "Pulling and starting services..."
cd /opt/crab-trap
docker-compose pull
docker-compose up -d

echo "Crab Trap started successfully!"
echo "Server URL: https://${subdomain}.${domain_name}"
