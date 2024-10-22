FROM golang:1.23-alpine AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o discogs_service .

FROM alpine:latest

RUN apk --no-cache add postgresql-client

WORKDIR /app

COPY --from=builder /app/discogs_service .

EXPOSE 8080

CMD ["./discogs_service"]
