# Use an official lightweight Go image
FROM golang:1.23.2-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o app .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./app"]
