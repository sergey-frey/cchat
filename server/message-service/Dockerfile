FROM golang:1.24.9-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGOENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/message-service ./cmd/message/main.go

RUN chmod +x ./wait-for-db.sh

FROM alpine:latest

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/message-service .
COPY --from=builder /src/config /config
COPY --from=builder /src/wait-for-db.sh .

USER appuser

EXPOSE 3030

CMD ["./message-service"]