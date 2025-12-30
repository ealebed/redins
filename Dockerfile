# syntax = docker/dockerfile:experimental
#
# Builder
#

FROM golang:1.25-alpine AS builder

# set working directorydoc
RUN mkdir -p /go/src/redins
WORKDIR /go/src/redins

# load dependency
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/mod go mod download

# copy sources
COPY . .

# build
RUN make

#
# ------ get latest CA certificates
#
FROM alpine:3.23 as certs
RUN apk --update add ca-certificates

#
# Runtime
#
FROM scratch

# copy CA certificates
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# this is the last command since it's never cached
COPY --from=builder /go/src/redins/.bin/github.com/ealebed/redins/redins /redins

CMD ["/redins"]
