FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.work go.work.sum ./
COPY pkg/ pkg/
COPY proto/ proto/
COPY services/user-service/ services/user-service/

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./services/user-service/cmd/

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/user-service .
COPY services/user-service/config/ config/

EXPOSE 50051

CMD ["./user-service"]