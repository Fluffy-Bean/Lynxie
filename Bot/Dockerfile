# syntax=docker/dockerfile:1
FROM alpine:3.18.2

RUN mkdir /lynxie && mkdir /lynxie/data
WORKDIR /lynxie

COPY ./lynxie /lynxie
COPY ./poetry.lock /lynxie
COPY ./pyproject.toml /lynxie

RUN apk update
RUN apk --no-cache add python3 py3-pip curl

RUN curl -sSL https://install.python-poetry.org | python3 -
RUN /root/.local/bin/poetry install
RUN /root/.local/bin/poetry run python3 database.py

CMD ["/root/.local/bin/poetry", "run", "python3", "main.py"]
