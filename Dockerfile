FROM golang:1.19-buster as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o crypto-bot .

FROM alpine:latest
RUN apk add tzdata
COPY --from=builder /app/crypto-bot .
CMD ["./crypto-bot"]