FROM golang:1.22-bookworm AS builder

WORKDIR /build
COPY app/ .

RUN go build -o authserver .

FROM debian:bookworm-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

COPY --from=builder /build/authserver /usr/local/bin/authserver

WORKDIR /app

COPY . .
RUN chmod +x test.sh

# Install Go toolchain in final image so agent can rebuild
RUN apt-get update -qq && apt-get install -y --no-install-recommends golang-go && \
    rm -rf /var/lib/apt/lists/*

CMD ["authserver"]
