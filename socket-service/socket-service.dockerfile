FROM alpine:latest

RUN mkdir /app

COPY sockApp /app
EXPOSE 3000
CMD [ "/app/sockApp"]