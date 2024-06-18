# Builder
FROM golang:1.20-alpine3.19 as builder

RUN apk update && apk add envsubst

RUN mkdir -p /opt/service

WORKDIR /opt/service

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN export CALENDAR_API_HTTP_HOST="0.0.0.0" && \
    export CALENDAR_API_HTTP_PORT="8080" && \
    export CALENDAR_API_GRPC_HOST="0.0.0.0" && \
    export CALENDAR_API_GRPC_PORT="50051" && \
    export CALENDAR_DB_DSN="postgres://calendar:calendar@db:5432/calendar" && \
    export CALENDAR_RABBITMQ_DSN="amqp://guest:guest@rabbitmq:5672" && \
    envsubst < ./configs/config.toml.dist > ./configs/calendar_config.toml

RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./bin/calendar ./cmd/calendar

# Release
FROM alpine:3.19.1

WORKDIR /root/
COPY --from=builder /opt/service/bin/calendar ./
COPY --from=builder /opt/service/configs/calendar_config.toml ./calendar_config.toml

CMD ["/root/calendar", "--config", "/root/calendar_config.toml"]