FROM debian:bookworm-slim

ENV JAVA_TOOL_OPTIONS="-Dfile.encoding=UTF-8"

RUN apt-get update -qq && apt-get install -y --no-install-recommends openjdk-17-jdk bash git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/ /app/app/
COPY test.sh /app/test.sh
RUN chmod +x test.sh

# Compile all Java sources
RUN mkdir -p /app/build && javac -d /app/build app/*.java

CMD ["./test.sh"]
