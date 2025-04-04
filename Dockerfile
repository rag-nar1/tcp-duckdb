FROM golang:1.24

WORKDIR /app

# Install SQLite tools
RUN apt-get update && apt-get install -y sqlite3 && apt-get clean

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /app/build/server main/*

# Create directories for storage
RUN mkdir -p /app/storge/server

# Create SQLite database file and initialize with schema
RUN touch /app/storge/server/db.sqlite3 && \
    sqlite3 /app/storge/server/db.sqlite3 < /app/storge/server/scheme.sql

# Set environment variables
ENV ServerPort=4000
ENV ServerAddr=0.0.0.0
ENV DBdir=/app/storge/
ENV ServerDbFile=server/db.sqlite3
ENV ENCRYPTION_KEY=A15pG0m3hwf0tfpVW6m92eZ6vRmAQA3C

# Expose the server port
EXPOSE 4000

# Run the server
CMD ["/app/build/server"] 