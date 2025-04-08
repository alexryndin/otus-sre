# Stage 1: Build the Go binary
FROM golang:1.22 AS go-builder

WORKDIR /app

# Copy the Go source file
COPY main.go .

# Build the Go binary (assuming no external dependencies)
RUN GO111MODULE=off go build -o goapp main.go

# Stage 2: Build the final image with Python and Go
FROM python:3.11-slim

WORKDIR /app

# Install Python dependencies
COPY app.py .
RUN pip install --no-cache-dir flask prometheus-flask-exporter api prometheus_client

# Copy the Go binary from the builder stage
COPY --from=go-builder /app/goapp .

# Expose ports for both apps
EXPOSE 3000 8080

# Start both apps
CMD ["sh", "-c", "python app.py & ./goapp"]
