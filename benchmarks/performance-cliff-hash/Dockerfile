FROM debian:bookworm-slim AS builder

ENV JAVA_TOOL_OPTIONS="-Dfile.encoding=UTF-8"

RUN apt-get update -qq && apt-get install -y --no-install-recommends openjdk-17-jdk

WORKDIR /build
COPY src/ .

RUN javac -d /build/classes *.java

FROM debian:bookworm-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash git openjdk-17-jre openjdk-17-jdk && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

COPY --from=builder /build/classes/ /app/classes/

WORKDIR /app

COPY . .
RUN chmod +x test.sh

# Create a wrapper script for running the Java program
RUN printf '#!/bin/sh\nexec java -cp /app/classes Main "$@"\n' > /usr/local/bin/perftest && \
    chmod +x /usr/local/bin/perftest

CMD ["perftest"]
