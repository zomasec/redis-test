# Use Go base image
FROM golang:1.19

RUN apt-get update && apt-get install -y \
    build-essential \
    telnet \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /go/src/app

# Copy Go modules and install dependencies
COPY go.mod .
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o scanner .

# Run the Go application
CMD ["./scanner"]
