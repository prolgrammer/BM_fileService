FROM golang:1.23.1-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR etc/app

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

ENTRYPOINT ["./main"]