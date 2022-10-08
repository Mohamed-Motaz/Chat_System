# syntax=docker/dockerfile:1

FROM golang:alpine

RUN apk add build-base

ADD . /go/src/myapp
WORKDIR /go/src/myapp


RUN go mod download
RUN go build -o myserver

EXPOSE 5555

ENTRYPOINT ["./myserver"]

#docker build -t go_server_image .    
#docker run --name go_container go_server_image -p 5555:5555