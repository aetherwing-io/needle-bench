FROM python:3.12-alpine

RUN apk add --no-cache bash curl git && \
    git config --global user.email "bench@needle-bench.cc" && \
    git config --global user.name "needle-bench"

WORKDIR /app

COPY app/requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir -r /tmp/requirements.txt

COPY . .
RUN chmod +x test.sh

EXPOSE 5000

CMD ["python", "app/app.py"]
