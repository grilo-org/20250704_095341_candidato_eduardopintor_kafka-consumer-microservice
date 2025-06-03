FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN go build -o consumer ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/consumer .
CMD ["./consumer"]
