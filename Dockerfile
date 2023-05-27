# Start from a base Go image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download
RUN go mod tidy

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o server ./cmd

# Expose the desired port (replace 8080 with your application's port)
EXPOSE 8080

# Set the entry point for the container
CMD ["./server"]
