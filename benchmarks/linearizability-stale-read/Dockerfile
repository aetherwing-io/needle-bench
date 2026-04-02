FROM golang:1.22-bookworm AS builder

WORKDIR /build
COPY app/ .

RUN go build -o kv-cluster .

FROM debian:bookworm-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

COPY --from=builder /build/kv-cluster /usr/local/bin/kv-cluster

WORKDIR /app

COPY . .
RUN chmod +x test.sh && \
    git init && git add -A && git commit -m 'init'

# Install Go toolchain in final image so agent can rebuild
RUN apt-get update -qq && apt-get install -y --no-install-recommends golang-go && \
    rm -rf /var/lib/apt/lists/*

CMD ["kv-cluster"]
