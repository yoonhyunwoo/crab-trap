#!/bin/bash
set -e

echo "Initializing Moltbook Prompt Injector..."

# Update system
apt-get update -y
apt-get install -y docker.io curl git

# Start Docker
systemctl start docker
systemctl enable docker

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Create application directory
mkdir -p /opt/moltbook-injector
cd /opt/moltbook-injector

# Create config.yaml
cat > config.yaml <<'CONFIG_EOF'
server:
  port: 8080
  log_dir: "./logs"

worker:
  moltbook_api_key: "${moltbook_api_key}"
  server_url: "http://localhost:8080"
  submolt: "${moltbook_submolt}"
  interval_minutes: ${worker_interval}
  os_detection: true

logging:
  level: "info"
  save_requests: true
CONFIG_EOF

# Create Dockerfile
cat > Dockerfile <<'DOCKERFILE_EOF'
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o worker ./cmd/worker

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/worker .
COPY --from=builder /app/prompts ./prompts

RUN mkdir -p /logs

EXPOSE 8080

CMD ["./server"]
DOCKERFILE_EOF

# Create docker-compose.yml
cat > docker-compose.yml <<'COMPOSE_EOF'
version: '3.8'

services:
  server:
    build: .
    command: ["./server"]
    ports:
      - "8080:8080"
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
        awslogs-stream-prefix: server

  worker:
    build: .
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
        awslogs-stream-prefix: worker

volumes:
  logs:
COMPOSE_EOF

# Create log directory
mkdir -p /opt/moltbook-injector/logs

echo "Building and starting services..."
cd /opt/moltbook-injector
docker-compose build
docker-compose up -d

echo "Moltbook Prompt Injector started successfully!"
echo "Server URL: http://localhost:8080"
