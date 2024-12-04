FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/rarimo/web3-auth-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/web3-auth-svc /go/src/github.com/rarimo/web3-auth-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/web3-auth-svc /usr/local/bin/web3-auth-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["web3-auth-svc"]
