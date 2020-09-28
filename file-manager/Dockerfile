# first stage
# how to copy dependancies??

# FROM golang:1.12.0-alpine3.9 AS builder
# RUN apk update && \
#     apk add --no-cache build-base imagemagick-dev git && \
#     mkdir /gopath
#
# COPY ${GOPATH}/src/ $GOPATH/src/
#
# ADD cmd cmd
# ADD configs configs
# ADD pkg pkg
#
# RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o file_manager ./cmd/file_manager

# Second stage

# FROM alpine:3.8
FROM centos:7

WORKDIR /opt

# COPY --from=builder /go/src/app/file_manager .

RUN yum install -y wget gcc make libpng-devel libjpeg-devel libungif-devel \
    yum clean all && \
    wget https://www.imagemagick.org/download/ImageMagick.tar.gz && \
    tar -xvf ImageMagick.tar.gz && \
    cd ImageMagick-7* && \
    ./configure --disable-dependency-tracking && \
    make && \
    make install && \
    ldconfig /usr/local/lib && \
    cd .. && \
    rm -r ImageMagick-7* ImageMagick.tar.gz && \
    yum remove -y wget gcc make libpng-devel && \
    yum clean all

# RUN apk update && apk add --no-cache imagemagick imagemagick-dev


ADD ./build ./

ENTRYPOINT ["./file_manager"]
