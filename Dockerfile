# -------- Stage 1: Builder --------
FROM golang:1.22 AS builder

WORKDIR /app

# Copy only go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy full source code
COPY . .

# Build binary (static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# -------- Stage 2: Runtime --------
FROM alpine:latest

# Add CA certs for outbound HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy built binary from builder
COPY --from=builder /app/main .

# Expose port (optional, just for documentation)
EXPOSE 8080

# Run the app
CMD ["./main"]
