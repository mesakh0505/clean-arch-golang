# Please keep up to date with the new-version of Golang docker for builder
FROM golang:1.14.0-stretch as builder

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y git \
    make openssh-client

WORKDIR /go_modules/LieAlbertTriAdrian/todo-service

# This is for private library that may be used in the projects
ENV GOPRIVATE="github.com/LieAlbertTriAdrian"

# Use git with SSH instead of https
RUN git config --global url."git@github.com:".insteadOf "https://github.com"
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

# Copy SSH key for git private repos
ADD .ssh/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

COPY . .
RUN make engine

## Distribution
FROM debian:stretch-slim

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y ca-certificates tzdata && \
    mkdir /app && mkdir todo-service

WORKDIR /todo-service

EXPOSE 9090

COPY --from=builder /go_modules/LieAlbertTriAdrian/todo-service/engine /app

CMD /app/engine rest