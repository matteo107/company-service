FROM golang:1.19.0 as builder

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy
RUN env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o=/companyservice ./cmd/api

FROM alpine:latest as production

WORKDIR /

COPY --from=builder /companyservice /companyservice

ENTRYPOINT ["./companyservice" ]