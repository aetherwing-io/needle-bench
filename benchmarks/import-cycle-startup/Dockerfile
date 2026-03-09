FROM python:3.12-alpine

RUN apk add --no-cache bash git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/ /app/app/
COPY test.sh /app/test.sh
RUN chmod +x test.sh

CMD ["./test.sh"]
