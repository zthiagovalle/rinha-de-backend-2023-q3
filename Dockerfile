# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rinha-api .

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/rinha-api .

EXPOSE 8080

CMD ["./rinha-api"]