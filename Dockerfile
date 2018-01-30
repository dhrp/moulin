FROM golang:1.8

RUN apt-get update && \
	apt-get install unzip &&\
	rm -rf /var/lib/apt/lists/

# install go dependencies
RUN go get github.com/mediocregopher/radix.v2


# install protoc
RUN curl -L https://github.com/google/protobuf/releases/download/v3.2.0/protoc-3.2.0-linux-x86_64.zip > protoc.zip && \
	unzip protoc.zip -d .
RUN go get google.golang.org/grpc

# install and build go app
COPY . /go/src/github.com/dhrp/moulin/
WORKDIR /go/src/github.com/dhrp/moulin/

RUN make -C protobuf
RUN make build

EXPOSE 50051
ENTRYPOINT ["/go/bin/server"]
