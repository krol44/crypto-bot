FROM golang:1.16-buster as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o crypto-bot .

FROM alpine:latest
MAINTAINER Krol44 <krol44@me.com>
COPY --from=builder /app/crypto-bot .
ENTRYPOINT ["./crypto-bot"]