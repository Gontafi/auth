FROM golang:1.21.4-alpine3.18 AS builder
WORKDIR /src
COPY . .
RUN go build -v -ldflags="-w -s" -o server main.go

FROM golang:1.21.4-alpine3.18
WORKDIR /app
COPY --from=builder /src/server .
COPY --from=builder /src/.env .