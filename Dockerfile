FROM golang:1.20

RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o trackme ./cmd/server

CMD ["./trackme"]
