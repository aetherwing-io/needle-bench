FROM golang:1.24-alpine

RUN apk add --no-cache git bash

WORKDIR /workspace

# Minimal reproduction of the k8s assume_cache silent drop bug
COPY app/ /workspace/

# The test verifies the fix
COPY test.sh /workspace/test.sh
RUN chmod +x /workspace/test.sh
