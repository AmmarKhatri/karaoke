FROM alpine:latest

RUN mkdir /app

COPY gameApp /app
COPY song.txt /app

CMD [ "/app/gameApp"]