FROM python:3.12.1-alpine AS base

RUN apk add --no-cache usbutils

WORKDIR /app

COPY ./requirements.txt .

RUN pip install --upgrade pip

RUN pip install --no-cache-dir -r requirements.txt

WORKDIR /app/src

COPY ./src .

# development stage
FROM base AS development

WORKDIR /app

COPY ./pyproject.toml .

COPY ./requirements-dev.txt .

RUN pip install --no-cache-dir -r requirements-dev.txt

#isort stage
FROM development AS isort

CMD [ "isort","./src/" ]

#black stage
FROM development AS black

CMD [ "black","./src/" ]

#mypy stage
FROM development AS mypy

CMD [ "mypy","./src/" ]

#pylint stage
FROM development AS pylint

WORKDIR /app

CMD [ "pylint","./src/" ]
