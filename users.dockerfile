FROM alpine:latest

RUN mkdir /app

COPY usersApp /app

CMD [ "/app/usersApp"]