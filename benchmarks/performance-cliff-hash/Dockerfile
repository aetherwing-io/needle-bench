FROM alpine:3.20 AS builder

RUN apk add --no-cache openjdk17-jdk

WORKDIR /build
COPY src/ .

RUN javac -d /build/classes *.java

FROM alpine:3.20

RUN apk add --no-cache bash coreutils git openjdk17-jre openjdk17-jdk && \
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
