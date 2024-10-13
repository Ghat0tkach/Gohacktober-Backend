# Use the official Go image with version 1.23
FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Copy the .env file
COPY .env ./

# Install dotenv (make sure to include it in your go.mod)
RUN go get github.com/joho/godotenv

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["go", "run", "cmd/server/main.go"]
