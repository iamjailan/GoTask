# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gotask ./cmd/api \
	&& CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
	&& addgroup -S gotask \
	&& adduser -S -G gotask gotask

WORKDIR /app

COPY --from=build /out/gotask /app/gotask
COPY --from=build /out/migrate /app/migrate
COPY migrations /app/migrations
COPY --chmod=755 docker-entrypoint.sh /app/docker-entrypoint.sh

USER gotask

ENV HTTP_ADDR=:8080 \
	MIGRATIONS_PATH=/app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
