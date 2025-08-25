# Docker Setup for Private GitHub Repository

This guide explains how to configure Docker builds to access the private healthcare-logging library.

## Problem

Docker builds fail when trying to download the private `github.com/louhibi/healthcare-logging` repository because:
1. Alpine Linux images don't include git by default
2. Private repositories require authentication

## Solution Options

### Option 1: Make Repository Public (Recommended)

The simplest solution for open source projects:

1. Go to https://github.com/louhibi/healthcare-logging
2. Settings → General → Danger Zone → Change visibility → Make public
3. No other changes needed - Docker builds will work immediately

### Option 2: Use GitHub Personal Access Token (Private Repo)

If you want to keep the repository private:

#### Step 1: Create GitHub Personal Access Token

1. Go to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Select these scopes:
   - `repo` (Full control of private repositories)
4. Copy the token (you won't see it again)

#### Step 2: Create .env File

Add to your `.env` file in the project root:

```bash
# GitHub token for private repository access
GITHUB_TOKEN=your_github_token_here
```

#### Step 3: Update Remaining Dockerfiles

The API Gateway, User Service, and Patient Service are already updated. Update the remaining services:

**Location Service Dockerfile:**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for private repository access
RUN apk update && apk add --no-cache git

# Configure git for private repositories (if needed)
ARG GITHUB_TOKEN
RUN if [ -n "$GITHUB_TOKEN" ]; then \\
        git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; \\
    fi

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o location-service .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/location-service .

# Expose port
EXPOSE 8084

# Command to run
CMD ["./location-service"]
```

**Appointment Service Dockerfile:**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for private repository access
RUN apk update && apk add --no-cache git

# Configure git for private repositories (if needed)
ARG GITHUB_TOKEN
RUN if [ -n "$GITHUB_TOKEN" ]; then \\
        git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; \\
    fi

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o appointment-service .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/appointment-service .

# Expose port
EXPOSE 8083

# Command to run
CMD ["./appointment-service"]
```

#### Step 4: Build with Token

```bash
# Build all services
GITHUB_TOKEN=your_token_here docker compose build

# Or export it once
export GITHUB_TOKEN=your_token_here
docker compose build
```

### Option 3: SSH Keys (Alternative)

If you prefer SSH authentication:

1. Add your SSH public key to GitHub
2. Update Dockerfiles to use SSH:

```dockerfile
# Configure SSH instead of HTTPS
ARG GITHUB_SSH_KEY
RUN if [ -n "$GITHUB_SSH_KEY" ]; then \\
        mkdir -p ~/.ssh && \\
        echo "$GITHUB_SSH_KEY" > ~/.ssh/id_rsa && \\
        chmod 600 ~/.ssh/id_rsa && \\
        ssh-keyscan github.com >> ~/.ssh/known_hosts && \\
        git config --global url."git@github.com:".insteadOf "https://github.com/"; \\
    fi
```

## Security Considerations

### For Option 2 (GitHub Token):
- ⚠️  Never commit tokens to git
- ✅  Use environment variables or Docker secrets
- ✅  Limit token scope to minimum required
- ✅  Rotate tokens regularly

### For Production:
- Use Docker secrets or Kubernetes secrets
- Consider using a CI/CD pipeline with proper secret management
- Use multi-stage builds to avoid including tokens in final image

## Verification

After updating, verify the build works:

```bash
# Test single service build
docker compose build user-service

# Test all services
docker compose build

# Check logs for success
docker compose build 2>&1 | grep -i "download"
```

## Troubleshooting

### Error: "git: executable file not found"
- Solution: Dockerfiles need `RUN apk add git`

### Error: "repository does not exist or requires authentication"
- Solution: Check GITHUB_TOKEN is set and valid
- Verify token has correct permissions (repo scope)

### Error: "bad credentials"
- Solution: Regenerate GitHub token
- Check token hasn't expired

### Build works locally but fails in CI/CD
- Solution: Set GITHUB_TOKEN in CI/CD environment variables
- For GitHub Actions, use secrets

## Recommended Approach

**For Development**: Use Option 1 (make repo public) 
**For Production**: Use proper secret management with Option 2