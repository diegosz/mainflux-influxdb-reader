###
# Mainflux MongoDB Writer Dockerfile
###

FROM golang:alpine
MAINTAINER Mainflux

ENV MONGO_HOST mongo
ENV MONGO_PORT 27017

###
# Install
###

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mainflux/mainflux-mongodb-reader
RUN cd /go/src/github.com/mainflux/mainflux-mongodb-reader && go install

# Dockerize
ENV DOCKERIZE_VERSION v0.2.0
ADD https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz dockerize.tar.gz
RUN tar -C /usr/local/bin -xzf dockerize.tar.gz && rm -f dockerize.tar.gz


###
# Run main command with dockerize
###
CMD dockerize -wait tcp://$MONGO_HOST:$MONGO_PORT \
				-timeout 10s /go/bin/mainflux-mongodb-writer -m MONGO_HOST

