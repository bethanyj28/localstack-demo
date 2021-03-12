From golang:1.16-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o demo main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/demo .
CMD ["./demo"]
