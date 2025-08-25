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

# Copy module file and download deps (will create initial go.sum)
COPY go.mod go.sum ./
RUN go mod download
# Copy the rest of the source
COPY . .
# Ensure go.sum contains all required module checksums
RUN go mod tidy
# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o location-service

FROM alpine:3.19
# Install CA certs for HTTPS calls if ever needed
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/location-service /usr/local/bin/location-service
EXPOSE 8084
ENV PORT=8084
CMD ["/usr/local/bin/location-service"]
