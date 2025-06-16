# Use the official Go image (single stage)
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code (cmd directory)
COPY ./cmd ./cmd

RUN go build -o app ./cmd/src

# Optional: expose your app's port (e.g., 8080)
EXPOSE 8080

# Run the app
ENTRYPOINT ["./app"]
