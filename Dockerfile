FROM golang:1.23-alpine AS build

WORKDIR /app
COPY . .

RUN go build -o tinyping ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=build /app .

CMD ["./tinyping"]