# Use the official Golang image as the base
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency management
COPY go.mod go.sum ./

# Download Go modules
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o main ./server

# Expose port 8080 for the server
EXPOSE 8080

# Set the command to run the server
CMD ["go", "run", "server/main.go"]