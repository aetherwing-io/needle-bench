FROM python:3.12-alpine

RUN apk add --no-cache bash git tzdata && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

RUN pip install --no-cache-dir pytz

WORKDIR /app

COPY app/ /app/app/
COPY test.sh /app/test.sh
RUN chmod +x test.sh

CMD ["./test.sh"]
