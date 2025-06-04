FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api_mock main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api_mock .
CMD ["./api_mock"]
