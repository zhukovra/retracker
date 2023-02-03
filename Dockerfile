FROM golang:1.19 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w" -o app ./cmd/retracker

FROM alpine:latest
WORKDIR /root/
COPY --from=builder app ./
EXPOSE 8080
ENTRYPOINT ["./app"]