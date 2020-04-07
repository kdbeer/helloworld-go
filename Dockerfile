# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM ubuntu:18.04 as builder

RUN apt-get update

RUN apt-get upgrade -y

RUN apt-get autoremove -y

RUN cat /etc/os-release

# Ubuntu
RUN apt-get install software-properties-common -y
RUN add-apt-repository ppa:kagamih/dlib -y
RUN apt-get install libdlib-dev libjpeg-turbo8-dev -y


RUN apt-get update

RUN cd /tmp

RUN apt-get install curl -y

ENV GOLANG_VERSION 1.4.2

RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz \
		| tar -v -C /usr/local -xz

ENV PATH /usr/local/go/bin:$PATH

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOROOT /usr/local/go

RUN mkdir -p /root/workspace/src/helloworld-go
ENV GOPATH /root/workspace
ENV PATH /go/bin:$PATH
WORKDIR /go

RUN go version
RUN go env

FROM golang:1.13 as builder2

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN pwd
RUN ls -alh

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server

CMD ["/server"]


# FROM alpine:3
# RUN apk add --no-cache ca-certificates

# # Copy the binary to the production image from the builder stage.
# COPY --from=builder /root/workspace/src/helloworld-go/* /

# RUN ls -alh /

# # Run the web service on container startup.
# CMD ["/server"]