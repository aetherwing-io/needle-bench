FROM rust:1.82-alpine AS builder

RUN apk add --no-cache musl-dev

WORKDIR /build
COPY Cargo.toml .
COPY src/ src/

RUN cargo build --release

FROM alpine:3.20

RUN apk add --no-cache bash coreutils python3 git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

COPY --from=builder /build/target/release/codegen-engine /usr/local/bin/codegen-engine

WORKDIR /app

COPY . .
RUN chmod +x test.sh

# Install Rust toolchain in final image so agent can rebuild
RUN apk add --no-cache curl gcc musl-dev && \
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain 1.82.0 && \
    rm -rf /var/cache/apk/*

ENV PATH="/root/.cargo/bin:${PATH}"

CMD ["codegen-engine"]
