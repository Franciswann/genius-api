# Stage 1: Build stage
# Use the official Golang image as the base environment
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files and download packages
# Doing this before copying the full source code allow Docker to cache layers
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Compile the application into a binary named "genius-app"
RUN go build -o genius-app .

# Stage 2: Final Stage
# Use a lightweight Alpine image for the runtime environment
FROM alpine:latest
WORKDIR /root/

# COPY the pre-compiled binary from the builder stage
COPY --from=builder /app/genius-app .

# Define the command to run the application
CMD ["./genius-app"]