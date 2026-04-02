FROM golang:1.21-bookworm

RUN apt-get update -qq && apt-get install -y --no-install-recommends bash curl git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/ ./
RUN go build -o server .

COPY test.sh .
RUN chmod +x test.sh

EXPOSE 8080

CMD ["./server"]
