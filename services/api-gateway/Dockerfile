FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.work go.work.sum ./
COPY pkg/ pkg/
COPY proto/ proto/
COPY services/api-gateway/ services/api-gateway/

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./services/api-gateway/cmd/

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/api-gateway .
COPY services/api-gateway/config/ config/

EXPOSE 8080

CMD ["./api-gateway"]