FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build
COPY app/ .

RUN go build -o server .

FROM alpine:3.20

RUN apk add --no-cache bash coreutils curl jq git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

# Install Go toolchain so agent can rebuild
RUN apk add --no-cache go gcc musl-dev

WORKDIR /app
COPY app/ app/
COPY --from=builder /build/server /app/server
COPY test.sh .
RUN chmod +x test.sh

CMD ["./server"]
