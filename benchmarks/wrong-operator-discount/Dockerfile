FROM python:3.12-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/ ./
COPY test.sh .
RUN chmod +x test.sh

CMD ["python", "main.py"]
