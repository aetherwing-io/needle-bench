FROM golang:1.21-bookworm AS builder

RUN apt-get update -qq && apt-get install -y --no-install-recommends gcc 

WORKDIR /build
COPY app/ .

RUN go build -o server .

FROM debian:bookworm-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash curl jq git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

# Install Go toolchain so agent can rebuild
RUN apt-get update -qq && apt-get install -y --no-install-recommends golang-go gcc 

WORKDIR /app
COPY app/ app/
COPY --from=builder /build/server /app/server
COPY test.sh .
RUN chmod +x test.sh

CMD ["./server"]
