# Start from the official Golang image
FROM golang:1.22-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git

# Copy go.mod and go.sum files to install dependencies
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o waitlist .

# Start a new minimal image for the application
FROM alpine:latest

# Set the working directory inside the new container
WORKDIR /app

# Install necessary packages in the minimal image
RUN apk add --no-cache ca-certificates

# Copy the built binary from the previous stage
COPY --from=builder /app/waitlist .

# Expose the port the app runs on
EXPOSE 3000

# Run the application
CMD ["./waitlist"]