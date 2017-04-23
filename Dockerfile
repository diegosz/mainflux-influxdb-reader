###
# Mainflux MongoDB Writer Dockerfile
###

FROM golang:alpine
MAINTAINER Mainflux

ENV INFLUX_HOST influx
ENV INFLUX_PORT 8086

###
# Install
###
# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mainflux/mainflux-influxdb-reader
RUN cd /go/src/github.com/mainflux/mainflux-influxdb-reader && go install

###
# Run main command with dockerize
###
CMD mainflux-influxdb-reader -i $INFLUX_HOST

