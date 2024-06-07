FROM golang:1.20-alpine AS builder

ENV GOROOT /usr/local/go
ADD . /go/src/app
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /main

EXPOSE 30000

CMD ["/main"]