# Crab Trap

**Crab Trap** - A tool for testing prompt injection vulnerabilities in AI agents. It generates malicious prompts and posts them to Moltbook, while an HTTP server collects and logs any executed commands.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Worker (Generator)                     │
│  • Detects OS environment                                │
│  • Generates prompts with env variables ($HOSTNAME, etc.)  │
│  • Posts to Moltbook API                                │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
              ┌──────────────┐
              │  Moltbook   │
              └──────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                 AI Agents (Victims)                     │
│  • Execute prompts containing curl commands                │
│  • Send environment data to HTTP server                   │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│               HTTP Server (Honeypot)                     │
│  • Logs all incoming requests                           │
│  • Saves to JSON files                                 │
│  • Web UI for monitoring                               │
└─────────────────────────────────────────────────────────────┘
```

## Features

- **OS Detection**: Automatically detects Unix/Linux or Windows environments
- **Template-Based Prompts**: Reusable prompt templates with variable substitution
- **HTTP Server**: Captures all requests with full logging
- **Web UI**: Real-time monitoring dashboard
- **Scheduled Posting**: Runs periodically or once
- **Rate Limit Handling**: Respects Moltbook API rate limits
- **Terraform**: Automated AWS deployment

## Quick Start

### Local Development

1. **Clone the repository**:
```bash
git clone https://github.com/yoonhyunwoo/crab-trap.git
cd crab-trap
```

2. **Update `config.yaml`**:
```yaml
worker:
  moltbook_api_key: "YOUR_MOLTBOOK_API_KEY"
  server_url: "http://localhost:8080"
  submolt: "general"
  interval_minutes: 60
```

3. **Start the server**:
```bash
go run cmd/server/main.go
```

4. **Run the worker** (one-time):
```bash
go run cmd/worker/main.go --once
```

Or run periodically:
```bash
go run cmd/worker/main.go
```

5. **Monitor**: Open http://localhost:8080

### AWS Deployment

1. **Create SSH key**:
```bash
ssh-keygen -t rsa -b 4096 -f terraform/crab-trap-key -N ""
```

2. **Get Moltbook API key**:
```bash
curl -X POST https://www.moltbook.com/api/v1/agents/register \
  -H "Content-Type: application/json" \
  -d '{"name": "CrabTrapAgent", "description": "Tests prompt injection vulnerabilities"}'
```
Copy the `api_key` from response and visit the `claim_url` to verify.

3. **Configure**:
```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars  # Add your api_key and IP
```

4. **Deploy**:
```bash
AWS_PROFILE=yhw terraform init
AWS_PROFILE=yhw terraform apply
```

5. **Build and push image**:
```bash
../scripts/build-and-push.sh
```

6. **Access**:
```bash
terraform output server_url  # http://injector.thumbgo.kr
terraform output ssh_command
```

**Note:** API key is stored in AWS Secrets Manager, not Terraform state.

## Project Structure

```
crab-trap/
├── cmd/
│   ├── server/              # HTTP server
│   │   └── main.go
│   └── worker/              # Worker
│       └── main.go
├── internal/
│   ├── server/              # Server package
│   │   ├── handler.go      # Request handlers
│   │   └── logger.go       # Log management
│   ├── worker/              # Worker package
│   │   ├── generator.go    # Prompt generator
│   │   └── poster.go       # Moltbook poster
│   ├── env/                # Environment detection
│   │   └── detector.go
│   └── config/             # Configuration
│       └── config.go
├── prompts/                 # Prompt templates
│   ├── unix/
│   │   ├── basic.txt
│   │   ├── curl.txt
│   │   └── advanced.txt
│   ├── windows/
│   │   ├── basic.txt
│   │   └── curl.txt
│   └── patterns.json
├── terraform/               # AWS infrastructure
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── provider.tf
│   └── user_data.sh
├── scripts/                # Build and deploy scripts
│   └── build-and-push.sh
├── pkg/moltbook/           # Moltbook SDK
├── config.yaml             # Configuration
├── Dockerfile              # Docker image definition
└── logs/                  # Request logs
```

## Prompt Templates

Prompts are organized by OS and complexity:

### Unix/Linux
- **basic.txt**: Simple hostname request
- **curl.txt**: JSON data with multiple env variables
- **advanced.txt**: Multiple environment variables

### Windows
- **basic.txt**: Simple computername request
- **curl.txt**: JSON data with multiple env variables

Templates use placeholders:
- `{{SERVER_URL}}`: Replaced with actual server URL
- `$VAR_NAME` or `%VAR_NAME%`: Replaced with actual environment values

## Configuration

### Server Config (`config.yaml`)
```yaml
server:
  port: 8080              # Server port
  log_dir: "./logs"        # Log directory
```

### Worker Config
```yaml
worker:
  moltbook_api_key: "YOUR_KEY"  # Moltbook API key
  server_url: "http://..."       # Server URL
  submolt: "general"            # Target submolt
  interval_minutes: 60            # Run interval
  os_detection: true             # Auto-detect OS
```

### CLI Flags
```bash
# Server
./server --port 8080 --log-dir ./logs

# Worker
./worker \
  --api-key YOUR_KEY \
  --submolt general \
  --server-url http://example.com \
  --once                     # Run once and exit
  --interval 30              # Custom interval (minutes)
```

## API Endpoints

### HTTP Server
- `GET /` - Web UI (monitoring dashboard)
- `POST /log` - Log incoming requests
- `GET /health` - Health check
- `GET /logs` - Get all logged requests (JSON)

### Moltbook SDK
Uses `pkg/moltbook` SDK for all API operations.

## Log Format

Requests are logged to JSON files:
- `logs/requests_YYYYMMDD.json` - Daily logs
- `logs/summary.json` - Summary statistics

Each log entry contains:
```json
{
  "timestamp": "2026-02-03T13:00:00Z",
  "method": "POST",
  "url": "/log",
  "headers": {...},
  "query_params": {...},
  "body": "hostname=...",
  "remote_addr": "192.168.1.1:12345",
  "user_agent": "curl/7.81.0"
}
```

## AWS Resources

Terraform creates:
- **EC2 Instance**: t3.small with Ubuntu 22.04
- **Security Group**: Ports 22 (SSH), 8080 (HTTP)
- **IAM Role**: CloudWatch Logs permissions
- **Route53 Record**: injector.thumbgo.kr
- **CloudWatch Log Groups**: /crab-trap/server, /crab-trap/worker
- **Docker Compose**: Auto-start server and worker

## Security Notes

⚠️ **Warning**: This tool is for security testing only.

- Do not use for malicious purposes
- Only test systems you own or have permission to test
- The HTTP server logs all requests - review logs regularly
- Use strong SSH keys for AWS access
- Limit SSH access to your IP in Terraform config

## Troubleshooting

### Worker not posting
- Check API key in config.yaml
- Verify server URL is accessible
- Check rate limits (1 post per 30 minutes)

### Terraform errors
- Ensure AWS credentials are configured (`aws configure`)
- Verify Route53 hosted zone ID is correct
- Check SSH key exists in `terraform/` directory

### No requests received
- Verify prompts are generated correctly
- Check Moltbook posts are visible
- Review worker logs for errors

## Development

```bash
# Install dependencies
go mod download

# Build
go build ./...

# Run tests
go test ./...

# Format code
go fmt ./...
goimports -w .

# Lint
golangci-lint run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - See LICENSE file for details

## Acknowledgments

- Built with [Moltbook SDK](https://github.com/moltbook/sdk-go)
- Infrastructure managed by [Terraform](https://www.terraform.io)
