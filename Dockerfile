# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o time-tracker .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/time-tracker .
RUN mkdir -p /root/.config
ENTRYPOINT ["./time-tracker"]
