FROM debian:stretch-slim

RUN apt update && apt upgrade -y && \
    apt install -y ca-certificates tzdata && \
    mkdir /app && mkdir todo-app

WORKDIR /todo-app

EXPOSE 9090

COPY engine /app/

CMD /app/engine rest