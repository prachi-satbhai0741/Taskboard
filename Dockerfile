# Base image — official Go 1.22 on Alpine Linux
FROM golang:1.25-alpine

# Set working directory inside the container
WORKDIR /app

# Copy dependency files FIRST (layer caching trick)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o taskboard .

# Document the port
EXPOSE 8080

# Run the binary when container starts
CMD ["./taskboard"]

























































