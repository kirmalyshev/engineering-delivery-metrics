FROM golang:1.24.4 as builder
WORKDIR /usr/local/go/src/delivery-dashboard
COPY go.mod ./
COPY go.sum ./
COPY Makefile ./
COPY pkg/ pkg/
COPY cmd/ cmd/
