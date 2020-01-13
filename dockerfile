FROM golang:1.13.6-alpine3.11
RUN apk add git gcc
COPY . /code
WORKDIR /code
RUN apk --update upgrade && \
    apk add sqlite && \
    rm -rf /var/cache/apk/*
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add libc-dev
RUN go get github.com/mattn/go-sqlite3
RUN go get -u github.com/gorilla/mux
RUN go mod init main
RUN go install
CMD main