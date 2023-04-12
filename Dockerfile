FROM golang:1.19.0

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy
RUN go build -o=/usr/local/bin/companyservice ./cmd/api

# FIXME: multi-stage build and copy binary from previous stage