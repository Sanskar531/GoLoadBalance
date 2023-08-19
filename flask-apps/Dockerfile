FROM python:slim-bullseye

RUN apt-get update && apt-get -y install build-essential
RUN pip install Flask uwsgi
WORKDIR /code/
COPY ./server.py .

CMD uwsgi --http :5000 --wsgi-file server.py --callable app -b 21573 --processes 5 --threads 5
