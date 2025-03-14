FROM golang:1.23 AS builder

WORKDIR /app

COPY .env ./
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o geo-service ./cmd/main.go

FROM alpine:3.18

WORKDIR /

COPY --from=builder /app/geo-service /geo-service
COPY --from=builder /app/.env /.env

RUN chmod +x /geo-service

RUN apk add --no-cache ca-certificates tzdata

CMD ["/geo-service"]