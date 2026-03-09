FROM alpine:3.20

RUN apk add --no-cache bash coreutils openjdk17-jdk git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app
COPY app/ app/

# Compile all Java sources
RUN javac -d classes app/*.java

COPY test.sh .
RUN chmod +x test.sh

CMD ["java", "-cp", "classes", "Main"]
