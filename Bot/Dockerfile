# syntax=docker/dockerfile:1
FROM alpine:3.18.2

EXPOSE 8000

RUN mkdir /app && mkdir /app/data
WORKDIR /app

COPY ./lynxie /app

RUN apk update && \
    apk --no-cache add python3 py3-pip && \
    pip install --no-cache-dir -r requirements.txt \

RUN python3 database.py

CMD ["python3", "main.py"]
