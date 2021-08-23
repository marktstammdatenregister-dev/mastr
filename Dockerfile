FROM python:3.8-slim

ENV PIP_DISABLE_PIP_VERSION_CHECK=on

RUN pip install poetry \
 && poetry config virtualenvs.create false
RUN apt-get -qq update \
 && apt-get -qq install gcc

WORKDIR /usr/src/app

COPY poetry.lock pyproject.toml .
RUN poetry install --no-interaction
