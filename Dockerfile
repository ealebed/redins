# syntax = docker/dockerfile:experimental
#
# Builder
#

FROM golang:1.27-alpine AS builder

# set working directorydoc
RUN mkdir -p /go/src/redins
WORKDIR /go/src/redins

# load dependency
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# copy sources
COPY . .

# build a static binary for scratch
RUN CGO_ENABLED=0 go build -o bin/redins ./

#
# ------ get latest CA certificates
#
FROM alpine:3.24 AS certs
RUN apk --update add ca-certificates

#
# Runtime
#
FROM scratch

# copy CA certificates
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# this is the last command since it's never cached
COPY --from=builder /go/src/redins/bin/redins /redins

CMD ["/redins"]
