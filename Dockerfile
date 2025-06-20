# Build stage
FROM golang:latest AS builder

WORKDIR /app

# Copy go.mod first to leverage Docker cache
COPY app/go.mod ./

# Download dependencies (this will create go.sum)
RUN go mod download

# Copy the rest of the source code
COPY app/ ./

# Tidy up dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run stage
FROM alpine:latest
# Install timezone data
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/main /app/main

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/main"]