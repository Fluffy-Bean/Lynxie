# syntax=docker/dockerfile:1
FROM alpine:3.18.2

RUN mkdir /app && mkdir /app/data
WORKDIR /app

COPY ./lynxie /app
COPY ./poetry.lock /app
COPY ./pyproject.toml /app

RUN apk update
RUN apk --no-cache add python3 py3-pip curl

RUN curl -sSL https://install.python-poetry.org | python3 -
RUN poetry install
RUN poetry run python3 database.py

CMD ["poetry", "run", "python3", "main.py"]
