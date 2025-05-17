FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o weatherapi .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/weatherapi .

COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./weatherapi"]
