# Please keep up to date with the new-version of Golang docker for builder
FROM golang:1.14.0-stretch

RUN apt update && apt upgrade -y && \
    apt install -y git \
    make openssh-client

WORKDIR /go_modules/LieAlbertTriAdrian/todo-service

# This is for private library that may be used in the projects
ENV GOPRIVATE="github.com/LieAlbertTriAdrian"

# Use git with SSH instead of https
RUN git config --global url."git@github.com:".insteadOf "https://github.com"
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air
