FROM golang:1.8.3
MAINTAINER harsh@portworx.com

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega

WORKDIR /
COPY . /

ENTRYPOINT ["ginkgo -v"]
CMD []
