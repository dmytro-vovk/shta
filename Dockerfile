# Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Souces
COPY . .

# Run build
RUN go build -o app cmd/main.go

# Runner
FROM alpine:latest
#RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Deploy the binary
COPY --from=builder /app/app .

# Make HTTP port externally available
EXPOSE 8080

# Run
CMD ["./app"]
