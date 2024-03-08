FROM golang:1.21.4-alpine

# Set destination for COPY
WORKDIR /app

# Create necessary directories
RUN mkdir -p /app/cmd/api /app/internal /app/sql

# Copy Go mod and sum files
COPY go.mod go.sum  ./

# Download Go modules
RUN go mod tidy

# Copy the source code
COPY . .


# Build
RUN go build -o gomovies ./cmd/api

# Expose port
EXPOSE 8080

# Run
CMD ["./gomovies"]
