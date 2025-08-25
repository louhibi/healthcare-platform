# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for private repository access
RUN apk update && apk add --no-cache git

# Configure git for private repositories (if needed)
ARG GITHUB_TOKEN
RUN if [ -n "$GITHUB_TOKEN" ]; then \
        git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; \
    fi

# Configure Go for private repositories
ENV GOPRIVATE=github.com/louhibi/*
ENV GONOPROXY=github.com/louhibi/*
ENV GONOSUMDB=github.com/louhibi/*

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Update go.sum with new dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o patient-service .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/patient-service .

# Expose port
EXPOSE 8082

# Command to run
CMD ["./patient-service"]