FROM node:20-slim

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash curl jq git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app
COPY app/ app/

# Install dependencies and compile TypeScript
RUN cd app && npm install && npx tsc

COPY test.sh .
RUN chmod +x test.sh

CMD ["node", "app/dist/server.js"]
