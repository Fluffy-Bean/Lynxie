# syntax=docker/dockerfile:1
FROM alpine:3.18.2

RUN mkdir /app && mkdir /app/data
WORKDIR /app

COPY ./lynxie /app
COPY ./requirements.txt /app

RUN apk update
RUN apk --no-cache add python3 py3-pip

RUN pip install --no-cache-dir -r requirements.txt
RUN python3 database.py

CMD ["python3", "main.py"]
