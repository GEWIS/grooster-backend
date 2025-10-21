# Use the official Go image (single stage)
FROM golang:1.24-alpine AS base

# Set the working directory inside the container
WORKDIR /app

RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code (cmd directory)
COPY ./cmd ./cmd

RUN CGO_ENABLED=1 go build -o app ./cmd/src

FROM alpine:3.22

WORKDIR /app

COPY --from=base /app/app ./app

# Run the app
ENTRYPOINT ["./app"]
