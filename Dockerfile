# Multistage docker image

## Build go binary for alpine 
FROM golang:1.11.9-alpine3.9 AS builder
RUN apk update && apk add --no-cache git
ENV GO111MODULE=on
WORKDIR /go/src/github.com/rraks/remocc
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./ 

## Build image from scratch
FROM alpine
### Install openssh and configure it
RUN apk --update add --no-cache openssh bash \
  && sed -i s/#PermitRootLogin.*/PermitRootLogin\ yes/ /etc/ssh/sshd_config \
  && echo "root:root" | chpasswd \
  && rm -rf /var/cache/apk/*
RUN /usr/bin/ssh-keygen -A
RUN ssh-keygen -t rsa -b 4096 -f  /etc/ssh/ssh_host_key
RUN mkdir /root/.ssh && touch /root/.ssh/authorized_keys
COPY ./setup/sshd/sshd_config ./root/.ssh/sshd_config
### Copy built image to final image
COPY --from=builder /go/src/github.com/rraks/remocc/ /go/src/github.com/rraks/remocc/
WORKDIR /go/src/github.com/rraks/remocc










