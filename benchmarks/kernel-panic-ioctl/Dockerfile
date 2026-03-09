FROM alpine:3.20 AS builder

RUN apk add --no-cache gcc musl-dev make

WORKDIR /build
COPY src/ .

RUN make

FROM alpine:3.20

RUN apk add --no-cache bash coreutils git gcc musl-dev make && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

COPY --from=builder /build/devctl /usr/local/bin/devctl

WORKDIR /app

COPY . .
RUN chmod +x test.sh

CMD ["devctl"]
