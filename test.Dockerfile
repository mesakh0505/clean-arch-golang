# Please keep up to date with the new-version of Golang docker for builder
FROM golang:1.13.1-stretch

RUN apt update && apt upgrade -y && \
    apt install git make openssh-client

WORKDIR /go_modules/LieAlbertTriAdrian/todo-service

# This is for private library that may be used in the projects
ENV GOPRIVATE="github.com/LieAlbertTriAdrian"

# Use git with SSH instead of https
RUN git config --global url."git@github.com:".insteadOf "https://github.com"
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

# Copy SSH key for git private repos
ADD .ssh/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

# This Docker will only to run the test
CMD make full-test
