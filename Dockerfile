FROM alpine:edge as builder
WORKDIR /app
RUN apk add --update go=1.16.7-r0 gcc g++
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o crypto-bot .

FROM golang:alpine
MAINTAINER Krol44 <krol44@me.com>
RUN apk add tzdata
COPY --from=builder /app/crypto-bot .
EXPOSE 3434
ENTRYPOINT ["./crypto-bot"]