FROM golang

COPY . /go/src/github.com/takkyuuplayer/go-anki
WORKDIR /go/src/github.com/takkyuuplayer/go-anki

RUN make
RUN make heroku

CMD ./main --addr ":${PORT}"
