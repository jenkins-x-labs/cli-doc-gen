FROM golang:1.12

WORKDIR /go/src/github.com/jstrachan/cli-doc-gen

COPY . /go/src/github.com/jstrachan/cli-doc-gen

RUN make linux

FROM centos:7

RUN yum install -y git

ENTRYPOINT ["cli-doc-gen"]

COPY --from=0 /go/src/github.com/jstrachan/cli-doc-gen/build/linux/cli-doc-gen /usr/local/bin
