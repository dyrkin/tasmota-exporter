FROM golang:1.18.4 AS builder

ARG ARCH=amd64

WORKDIR /tasmota-exporter/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH go build -o app ./cmd

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /tasmota-exporter/app app

CMD ["./app"]