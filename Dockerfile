FROM golang:1.21.0-alpine3.18 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app/service

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN go build -o ./bin/service ./cmd/main.go

FROM alpine:3.18

WORKDIR /app/bin

COPY --from=builder /app/service/bin .

ENTRYPOINT [ "./service" ]