FROM golang:alpine AS builder

WORKDIR /go/src/github.com/Strum355/2Cloud

ENV GOBIN=/go/2Bot
ENV GOPATH=/go

COPY --from=builder /go/2Bot /go/2Bot
