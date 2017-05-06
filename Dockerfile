FROM golang:alpine
LABEL Description="Docker image of the Nyaa replacement"

# Get build reqs
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev

# Get build deps
RUN go get github.com/gorilla/feeds && \
    go get github.com/gorilla/mux && \
    go get github.com/mattn/go-sqlite3 && \
    go get github.com/jinzhu/gorm && \
    go get github.com/Sirupsen/logrus && \
    go get gopkg.in/natefinch/lumberjack.v2

# Build
RUN mkdir -p /go/src/github.com/ewhal/nyaa
WORKDIR /go/src/github.com/ewhal/nyaa
ADD . /go/src/github.com/ewhal/nyaa/
RUN go build -o nyaa .

ENTRYPOINT ["./nyaa", "-host", "0.0.0.0"]