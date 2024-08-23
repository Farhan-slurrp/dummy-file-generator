FROM golang:1.21.6 as base

WORKDIR /app

COPY templates templates
COPY go.mod go.mod
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-mod=vendor go build -o main ./main.go

FROM alpine:latest as final

WORKDIR /

COPY --from=base /app/main .

EXPOSE 8080

ENTRYPOINT ["./main"]