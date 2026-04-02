FROM node:20-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash curl git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/ ./
COPY test.sh .
RUN chmod +x test.sh

EXPOSE 3000

CMD ["node", "server.js"]
