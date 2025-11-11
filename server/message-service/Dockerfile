FROM golang:1.24.9-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./wait-for-db.sh ./

RUN CGOENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/messages-service ./cmd/message/main.go

RUN chmod +x ./wait-for-db.sh

FROM alpine:latest

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt ./etc/ssl/certs/
COPY --from=builder /app/messages-service .
COPY --from=builder /src/config ./config
COPY --from=builder /src/wait-for-db.sh .
COPY --from=builder /src/.env .
COPY --from=builder /src/schema ./schema

RUN chmod +x ./wait-for-db.sh

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

EXPOSE 8080

CMD ["./messages-service"]