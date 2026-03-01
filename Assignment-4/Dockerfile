# Stage 1 - build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/api/main.go

# Stage 2 - minimal image
FROM alpine
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8080
ENTRYPOINT ["./app"]