FROM alpine:3.19

RUN apk add --no-cache openjdk17-jdk bash git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/ /app/app/
COPY test.sh /app/test.sh
RUN chmod +x test.sh

# Compile all Java sources
RUN mkdir -p /app/build && javac -d /app/build app/*.java

CMD ["./test.sh"]
