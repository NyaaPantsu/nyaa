FROM golang:alpine
LABEL Description="Docker image of the Nyaa replacement"

# Get build reqs
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev

# Build
RUN mkdir -p /go/src/github.com/ewhal/nyaa
WORKDIR /go/src/github.com/ewhal/nyaa
ADD . /go/src/github.com/ewhal/nyaa/
RUN go list -f '{{.Deps}}' | tr "[" " " | tr "]" " " | xargs go list -e -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | grep -v 'github.com/ewhal/nyaa' | xargs go get -v
RUN go build -o nyaa .

ENTRYPOINT ["./nyaa", "-host", "0.0.0.0"]