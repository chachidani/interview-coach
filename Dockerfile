# Use the official Golang image
FROM golang:1.22-alpine


# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -o app .

# Command to run the executable
CMD ["./app"]
