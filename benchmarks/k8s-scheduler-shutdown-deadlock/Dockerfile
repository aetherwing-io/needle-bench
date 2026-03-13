FROM golang:1.24-alpine

RUN apk add --no-cache git bash

WORKDIR /workspace

COPY app/ /workspace/
COPY test.sh /workspace/test.sh
RUN chmod +x /workspace/test.sh
