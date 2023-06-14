# Use the official Go image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and cache Go module dependencies
RUN go mod download

# Copy the application source code
COPY main.go config.go ./

# Build the Go application
RUN go build -o myapp

# Set the command to run the application
CMD ["./myapp"]
