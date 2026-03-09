FROM alpine:3.20

RUN apk add --no-cache bash coreutils curl jq python3 py3-pip sqlite git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app
COPY app/ app/
COPY test.sh .
RUN chmod +x test.sh

CMD ["python3", "app/server.py"]
