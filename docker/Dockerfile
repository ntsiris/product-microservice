# Stage 1: Build the Go application
FROM golang:1.22.7-alpine AS builder

WORKDIR /app
COPY . .

# Download dependencies and build the application
RUN go mod download && \
    go build -o product-microservice ./cmd

# Stage 2: Create the final image
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/product-microservice /app/product-microservice

# Expose the application port
EXPOSE 8080

# Set environment variables for production
ENV PORT=8080
ENV MIGRATE_UP=true
ENV MIGRATE_DOWN=false

# Run the application
CMD ["/app/product-microservice"]
