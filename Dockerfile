FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o taskmanager .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/taskmanager .

EXPOSE 8080
EXPOSE 50051

CMD ["./taskmanager"]