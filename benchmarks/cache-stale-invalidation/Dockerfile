FROM debian:bookworm-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash curl jq python3 git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app
COPY app/ app/
COPY test.sh .
RUN chmod +x test.sh

CMD ["python3", "app/server.py"]
