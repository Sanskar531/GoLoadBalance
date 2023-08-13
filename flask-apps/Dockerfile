FROM python:slim-bullseye

ENV FLASK_RUN_HOST=0.0.0.0

RUN pip install Flask
WORKDIR /code/
COPY ./server.py .

CMD flask --app server run
