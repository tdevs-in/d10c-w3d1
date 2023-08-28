# Use an official Go runtime as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code to the container
COPY . .

# Build the Go application
RUN go build -o todo-app .

# Expose the port that the application will listen on
EXPOSE 9999

# Command to run the Go application when the container starts
CMD ["./todo-app"]
