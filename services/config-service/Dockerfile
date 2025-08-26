FROM golang:1.21-alpine AS builder
ARG GITHUB_TOKEN
WORKDIR /app

# Install git (needed for private modules)
RUN apk add --no-cache git ca-certificates build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o config-service .

FROM alpine:3.19
WORKDIR /root/
RUN adduser -D appuser
USER appuser
COPY --from=builder /app/config-service ./
EXPOSE 8085
ENTRYPOINT ["./config-service"]
