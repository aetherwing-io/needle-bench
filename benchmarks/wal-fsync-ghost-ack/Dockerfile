FROM golang:1.22-bookworm AS builder

WORKDIR /build
COPY app/ .

RUN go build -o kvwal ./...

FROM debian:bookworm-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

# Install Go toolchain so the agent can rebuild after edits
RUN apt-get update -qq && apt-get install -y --no-install-recommends golang-go 

COPY --from=builder /build/kvwal /usr/local/bin/kvwal

WORKDIR /workspace
COPY app/ .

COPY test.sh .
RUN chmod +x test.sh

CMD ["kvwal", "serve"]
