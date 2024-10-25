# Use the official Go image but specify the correct platform
FROM --platform=linux/amd64 golang:1.23.2-alpine

# Set the working directory
WORKDIR /app

# Copy all files from the local directory to the container
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o app .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./app"]
