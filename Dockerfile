# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the gRPC server
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

# Use a minimal base image to reduce the final image size
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server .

# Expose the port on which the server will run
EXPOSE 50051
EXPOSE 2112

# Command to run the server when the container starts
CMD ["./server"]