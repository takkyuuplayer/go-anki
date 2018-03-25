FROM golang

RUN go get -u github.com/golang/dep/cmd/dep && \
    go get -u github.com/jessevdk/go-assets-builder

COPY . /go/src/github.com/takkyuuplayer/go-anki
WORKDIR /go/src/github.com/takkyuuplayer/go-anki

RUN make heroku

CMD ./main --addr ":${PORT}"
