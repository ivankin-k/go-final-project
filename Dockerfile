# ----- BUILD -----
FROM golang:1.26.3 AS builder
WORKDIR /build
COPY go.mod go.sum *.go ./
RUN go mod download
COPY pkg/ ./pkg/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# ----- RUN -----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/app ./run
COPY web/ ./web/
COPY dummy.db ./scheduler.db
CMD ["./run"]