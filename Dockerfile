FROM golang:1.6

COPY . /go/src/github.com/generationtux/brizo
RUN cd /go/src/github.com/generationtux/brizo && \
    export CGO_ENABLED=0 go build

WORKDIR /go/src/github.com/generationtux/brizo
ENTRYPOINT /go/src/github.com/generationtux/brizo/brizo
