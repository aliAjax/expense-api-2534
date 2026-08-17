# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/expense-api ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
RUN mkdir -p /app/data && chown -R app:app /app

COPY --from=build /out/expense-api /app/expense-api

USER app

EXPOSE 8080

ENV PORT=8080
ENV DB_PATH=/app/data/expense.db

CMD ["/app/expense-api"]
