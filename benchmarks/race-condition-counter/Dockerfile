FROM node:20-alpine

RUN apk add --no-cache redis curl bash git python3 && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/package.json /tmp/package.json
RUN cd /tmp && npm install --production && mv node_modules /app/node_modules

COPY . .
RUN chmod +x test.sh

EXPOSE 3000

CMD ["node", "app/server.js"]
