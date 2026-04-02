FROM debian:bookworm-slim

ENV JAVA_TOOL_OPTIONS="-Dfile.encoding=UTF-8"

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash openjdk-17-jdk git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app
COPY app/ app/

# Compile all Java sources
RUN mkdir -p classes && javac -d classes app/*.java

COPY test.sh .
RUN chmod +x test.sh

CMD ["java", "-cp", "classes", "Main"]
