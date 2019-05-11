FROM golang:1.12-alpine

ARG ostype=Linux

RUN apk --no-cache add \
    g++ \
    git \
    bash

ENV GOPROXY=https://gocenter.io

# Mock creator
RUN go get -u github.com/vektra/mockery/.../

# Create user
ARG uid=1000
ARG gid=1000

RUN bash -c 'if [ ${ostype} == Linux ]; then addgroup -g $gid app; else addgroup app; fi && \
    adduser -D -u $uid -G app app && \
    chown app:app -R /go'

# Fill go mod cache.
RUN mkdir /tmp/cache
COPY go.mod /tmp/cache
COPY go.sum /tmp/cache
RUN chown app:app -R /tmp/cache
USER app
RUN cd /tmp/cache && \
    go mod download 


WORKDIR /src
