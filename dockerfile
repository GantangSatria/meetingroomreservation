# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/app /app/app

EXPOSE 8080
CMD ["/app/app"]

