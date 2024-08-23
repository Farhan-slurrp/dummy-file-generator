FROM golang:1.21.6 as base

WORKDIR /dummy-file-generator

COPY templates templates
COPY vendor vendor
COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-mod=vendor go build -o main ./main.go

FROM alpine:latest as final

WORKDIR /

COPY --from=base /dummy-file-generator/main .

EXPOSE 8080

ENTRYPOINT ["./dummy-file-generator"]