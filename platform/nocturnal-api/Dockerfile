# Use the official Golang image as the base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory to install dependencies first
COPY nocturnal-api/go.mod nocturnal-api/go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire nocturnal-api directory contents (including subdirectories and files in /internal and /api)
COPY nocturnal-api/ .

# Expose the application port
EXPOSE 9000

# Build the application
RUN go build -o main ./cmd/app/main.go

# Specify the entry point for the container
CMD ["/app/main"]
