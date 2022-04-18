FROM golang:1.17.5-alpine

ENV GO111MODULE=on
ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org
ENV GOTESTSUM_FORMAT=testname
ARG version=develop
ARG debugBuild

RUN apk add --no-cache gcc libc-dev git
RUN go install gotest.tools/gotestsum@v1.7.0

WORKDIR /home/keptn-build-base
ADD go-dep.mod ./go.mod
RUN go mod download -x