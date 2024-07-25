FROM golang:1.21 AS builder

RUN apt-get update && apt-get install -y \
  gcc-aarch64-linux-gnu \
  libc6-dev-arm64-cross

WORKDIR /app

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64
ENV CC=aarch64-linux-gnu-gcc

RUN go build -o trackme ./cmd/server

FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/trackme .

CMD ["./trackme"]
