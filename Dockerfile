From golang:1.16-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o demo main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/demo .

# a bit of magic to wait for dependencies to be available for testing
ENV WAIT_VERSION 2.7.3
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

ENTRYPOINT /wait && ./demo
