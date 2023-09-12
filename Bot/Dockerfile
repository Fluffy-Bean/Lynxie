# syntax=docker/dockerfile:1
FROM alpine:3.18.2

# Make a directory for the app
RUN mkdir /app
RUN mkdir /app/lynxie
WORKDIR /app

# Copy the app files
COPY ./lynxie/ /app/lynxie/
COPY ./poetry.lock /app
COPY ./pyproject.toml /app

# Install dependencies
RUN apk update
RUN apk --no-cache add python3 py3-pip curl

# Install poetry
RUN curl -sSL https://install.python-poetry.org | python3 -
RUN export PATH="/root/.local/bin:$PATH"

# Install the app
RUN poetry install
RUN poetry run python3 database.py

# Run the app
CMD ["poetry", "run", "python3", "main.py"]
