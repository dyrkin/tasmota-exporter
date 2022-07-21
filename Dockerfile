FROM golang:1.18.4 AS builder
WORKDIR /tasmota-exporter/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /tasmota-exporter/app app

CMD ["./app"]