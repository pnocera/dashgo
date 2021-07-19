FROM golang:alpine

RUN mkdir -p /go/src/github.com/pnocera/dashgo
WORKDIR /go/src/github.com/pnocera/dashgo

build:
    COPY . /go/src/github.com/pnocera/dashgo/
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build \
        -ldflags='-w -s -extldflags "-static"' -a -installsuffix nocgo -o /server github.com/pnocera/dashgo/

    #ENTRYPOINT ["./server"]
    
