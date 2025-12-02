# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o time-tracker .

# Runtime stage
FROM alpine:latest
# Include nano as default editor fallback for the 'edit' command
RUN apk --no-cache add ca-certificates nano
WORKDIR /root/
COPY --from=builder /app/time-tracker .
RUN mkdir -p /root/.config
ENTRYPOINT ["./time-tracker"]
