FROM golang:1-alpine

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /go/src/dr-feeder
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT [ "/go/bin/dr-feeder", "-c", "/go/etc" ]
