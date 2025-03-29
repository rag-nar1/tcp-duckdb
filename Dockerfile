FROM golang:1.24

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /app/build/server main/*

# Create directory for storage
RUN mkdir -p /app/storge/server

# Expose the server port
EXPOSE 4000

# Run the server
CMD ["/app/build/server"] 